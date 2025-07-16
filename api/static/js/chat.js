// 多用户WebSocket连接管理
let websocketConnections = new Map(); // 存储每个用户的WebSocket连接
let reconnectTimers = new Map(); // 存储每个用户的重连定时器
let retryCounts = new Map(); // 存储每个用户的重试次数
// 最大重试次数（可根据需求调整）
const MAX_RETRIES = 10;
// 重连间隔（毫秒）
const RECONNECT_INTERVAL = 5000;
// 用于消息去重的缓存
let processedMessages = new Set();


// 获取当前用户账号的函数
function getCurrentUserAccount() {
    // 首先尝试从URL参数获取
    const urlParams = new URLSearchParams(window.location.search);
    const accountFromUrl = urlParams.get("account_num");
    if (accountFromUrl) {
        return accountFromUrl;
    }

    // 然后尝试从sessionStorage获取当前活跃用户
    const currentUser = sessionStorage.getItem("account_num");
    if (currentUser) {
        return currentUser;
    }

    // 最后尝试从localStorage的用户令牌中获取
    const userTokens = JSON.parse(localStorage.getItem('userTokens') || '{}');
    const userAccounts = Object.keys(userTokens);
    if (userAccounts.length > 0) {
        // 返回第一个用户的账号（或者可以根据登录时间选择最新的）
        return userAccounts[0];
    }

    return null;
}

// 获取当前用户的JWT token
function getCurrentUserJWT() {
    const currentUser = getCurrentUserAccount();
    if (!currentUser) {
        return null;
    }

    // 从多用户存储中获取当前用户的JWT
    const userTokens = JSON.parse(sessionStorage.getItem('userTokens') || '{}');
    const userToken = userTokens[currentUser];
    if (userToken && userToken.jwt) {
        return userToken.jwt;
    }

    // 如果没有找到，尝试从localStorage获取（向后兼容）
    const jwt =sessionStorage.getItem('jwt');
    if (jwt) {
        return jwt;
    }

    return null;
}

// 绑定发送消息的表单事件
document.getElementById("textInputForm").addEventListener("submit", function (event) {
    event.preventDefault();
    const input = document.getElementById("textInput");
    const message = input.value;

    if (!getCurrentFriend()) {
        alert("请先选择一个好友！");
        return;
    }

    const userInput = document.getElementById('textInput').value;
    if (!userInput.trim()) return; // 空输入不发送

    // 检查消息长度
    if (userInput.length > 1000) {
        alert("消息长度不能超过1000字符！");
        return;
    }

    const jsonData = JSON.stringify({
        data: userInput,
        accountNum: getCurrentUserAccount(),
        receiver: getCurrentFriend(),
        type: "SendText",
    });

    // 检查连接状态后再发送
    const currentWs = websocketConnections.get(getCurrentUserAccount());
    if (currentWs && currentWs.readyState === WebSocket.OPEN) {
        currentWs.send(jsonData);
        AddMessage(message, 'sent'); // 立即本地显示
        input.value = "";
        // 发送后重置输入框高度
        input.style.setProperty('height', '20px', 'important');
        // 发送后立即清空字符计数器
        const charCounter = document.getElementById('charCounter');
        if (charCounter) {
            charCounter.textContent = '0/1000';
            charCounter.classList.remove('warning', 'danger');
        }
        // 发送后清空当前好友的对话内容
        if (getCurrentFriend()) {
            clearFriendDraft(getCurrentFriend());
        }
        // 发送后自动滚动到最后一条消息
        scrollToBottom();
    } else {
        alert("连接未建立，无法发送消息");
        console.error("发送失败：WebSocket未连接");
    }
});

fetch(localStorage.getItem('base_url')+`/api/v1/user/info?account_num=${getCurrentUserAccount()}`, {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json',
        'Authorization': getCurrentUserJWT() ? `Bearer ${getCurrentUserJWT()}` : ''
    },
})
    .then(res => res.json())
    .then(data => {
        if (data.code !== 200) {
            document.getElementById('userNickname').textContent = '用户信息加载失败';
            document.getElementById('userAccountNum').textContent = '';
            return;
        }
        const user = data.data;
        document.getElementById('userNickname').textContent = user.name || '未设置昵称';
        document.getElementById('userAccountNum').textContent = '账号：' + (user.account_num || "");
        document.getElementById('userAvatar').src = getUserAvatar();
    })
    .catch(err => {
        document.getElementById('userNickname').textContent = '用户信息加载失败';
    });

document.addEventListener('DOMContentLoaded', function() {
    function setWebsocketConnections() {
        // 获取当前用户信息
        const currentUser = getCurrentUserAccount();
        if (!currentUser) {
            console.error('用户账号不存在，无法建立WebSocket连接');
            return;
        }

        // 清理其他用户的连接（如果有的话）
        for (const [userAccount, ws] of websocketConnections.entries()) {
            if (userAccount !== currentUser) {
                cleanupUserConnection(userAccount);
            }
        }

        let token;
        let url;
        // 直接使用WebSocket连接，不进行fetch预检查
        let jwtToken = getCurrentUserJWT();
        if (!jwtToken) {
            console.error('JWT token不存在，无法建立WebSocket连接');
            document.getElementById('error').textContent = 'JWT token不存在，请重新登录';
            return;
        }

        // 使用正确的端口9090
        token = jwtToken;
        url = `8.137.19.65:9090`; // 建议明确添加ws://协议前缀
        console.log("直接连接WebSocket:", { token: token ? '***' : 'null', url: url });
        websocketHandler(token, url);
    }

    // 执行函数
    setWebsocketConnections();
});

function websocketHandler(token, url) {
    const currentUser = getCurrentUserAccount();
    if (!currentUser) {
        console.error('用户账号不存在，无法建立WebSocket连接');
        return;
    }
    
    // 清除可能存在的旧定时器
    const oldTimer = reconnectTimers.get(currentUser);
    if (oldTimer) {
        clearTimeout(oldTimer);
        reconnectTimers.delete(currentUser);
    }
    // 重置重试次数
    retryCounts.set(currentUser, 0);

    // 创建WebSocket连接的内部函数
    function createWebSocket() {
        // 关闭可能存在的旧连接
        const oldWs = websocketConnections.get(currentUser);
        if (oldWs) {
            oldWs.close(1000, "重新连接");
            websocketConnections.delete(currentUser);
        }

        // 创建新连接
        const wsUrl = `ws://${url}/api/v1/chat/ws?token=${token}`;
        console.log(`尝试连接WebSocket: ${wsUrl}`);
        const ws = new WebSocket(wsUrl);
        
        // 存储当前用户的WebSocket连接
        websocketConnections.set(currentUser, ws);

        // 连接成功事件（保留原有逻辑）
        ws.onopen = function (event) {
            console.log('连接已打开', event);
            console.log('WebSocket连接成功，状态:', ws.readyState);
            AddMessage('已连接到服务器。', 'sent');
            // 连接成功后重置重试次数
            retryCounts.set(currentUser, 0);
        };

        // 接收消息事件（保留原有逻辑）
        ws.onmessage = function (event) {
            console.log('收到WebSocket消息:', event.data);
            const message = JSON.parse(event.data);
            console.log('解析后的消息:', message);
            console.log('消息类型:', message.type);
            console.log('消息字段:', Object.keys(message));
            switch (message.type) {
                case "RefreshText":
                    console.log('收到RefreshText消息');
                    break;
                case "SendText":
                    console.log('收到SendText消息，调用sendTextResponse');
                    sendTextResponse(message);
                    break;
                case "Login":
                    console.log('收到Login消息');
                    break;
                case "LogOff":
                    console.log('收到LogOff消息');
                    break;
                case "Register":
                    console.log('收到Register消息');
                    break;
                case "kick":
                    console.log('收到kick消息');
                    KickUser(message);
                    // 被踢出时停止重连
                    stopReconnect();
                    break;
                case "error":
                    console.log('收到error消息');
                    alert(message.data);
                    break;
                default:
                    console.log('收到未知类型的消息:', message.type);
            }
        };


        // 连接关闭事件（添加重连逻辑）
        ws.onclose = function (event) {
            console.log('连接已关闭', event);
            console.log('关闭详情:', {
                code: event.code,
                reason: event.reason,
                wasClean: event.wasClean
            });
            
            AddMessage('已从服务器断开。', 'sent');
            
            // 从连接映射中移除
            websocketConnections.delete(currentUser);
            
            // 正常关闭（1000）时不重连
            if (event.code === 1000) {
                console.log('正常关闭，不进行重连');
                return;
            }
            
            // 被踢出时不重连（检查reason是否存在）
            if (event.reason && (event.reason.includes("被踢出") || event.reason === "用户被踢出")) {
                console.log('被踢出，不进行重连');
                return;
            }

            // 其他情况触发重连
            console.log('异常关闭，准备重连');
            scheduleReconnect();
        };

        // 错误事件（添加重连逻辑）
        ws.onerror = function (error) {
            console.error('WebSocket 错误', error);
            console.error('错误详情:', {
                errorType: error.type,
                readyState: ws ? ws.readyState : 'no ws',
                url: ws ? ws.url : 'no url'
            });
            AddMessage('连接出错，正在尝试重连...', 'sent');
            // 错误时触发重连
            if (ws && ws.readyState !== WebSocket.CLOSED) {
                ws.close(1011, "连接错误");
            }
        };
    }

    // 安排重连的函数
    function scheduleReconnect() {
        const currentRetryCount = retryCounts.get(currentUser) || 0;
        
        // 检查是否超过最大重试次数
        if (currentRetryCount >= MAX_RETRIES) {
            console.error(`已达到最大重试次数（${MAX_RETRIES}次），停止重连`);
            AddMessage('连接服务器失败，请刷新页面重试', 'sent');
            return;
        }

        // 计算下次重连时间（可扩展为指数退避，这里用固定间隔）
        const newRetryCount = currentRetryCount + 1;
        retryCounts.set(currentUser, newRetryCount);
        const delay = RECONNECT_INTERVAL;
        console.log(`第${newRetryCount}次重连将在${delay/1000}秒后进行`);
        AddMessage(`连接断开，正在尝试第${newRetryCount}次重连...`, 'sent');

        // 设置重连定时器
        const timer = setTimeout(() => {
            createWebSocket();
        }, delay);
        reconnectTimers.set(currentUser, timer);
    }

    // 停止重连的函数（被踢时调用）
    function stopReconnect() {
        const timer = reconnectTimers.get(currentUser);
        if (timer) {
            clearTimeout(timer);
            reconnectTimers.delete(currentUser);
        }
    }

    // 初始化连接
    createWebSocket();
}


// 接收到sendText响应
function sendTextResponse(message) {
    console.log('处理SendText消息:', message);
    console.log('消息所有字段:', Object.keys(message));
    console.log('消息内容:', JSON.stringify(message, null, 2));
    
    // 获取当前用户账号
    const myAccount = getCurrentUserAccount();
    const chattingWith = getCurrentFriend();
    
    console.log('当前用户账号:', myAccount);
    console.log('当前聊天对象:', chattingWith);
    console.log('消息发送者:', message.accountNum);
    console.log('消息接收者:', message.receiver);

    // 获取消息的发送者和接收者，支持多种字段名格式
    const sender = message.accountNum || message.senderAccountNum || message.sender || '';
    const receiver = message.receiver || message.receiverAccountNum || '';
    
    // 检查消息是否属于当前聊天窗口
    if ((sender === chattingWith && receiver === myAccount) ||
        (sender === myAccount && receiver === chattingWith)) {
        
        // 检查消息是否已经处理过
        const messageKey = `${sender}-${receiver}-${message.data}`;
        if (processedMessages.has(messageKey)) {
            console.log('消息已处理过，跳过:', messageKey);
            return;
        }
        
        // 标记消息为已处理
        processedMessages.add(messageKey);
        
        // 限制缓存大小，避免内存泄漏
        if (processedMessages.size > 100) {
            const firstKey = processedMessages.values().next().value;
            processedMessages.delete(firstKey);
        }
        
        AddMessage(message.data, 'received');
        console.log('显示接收的消息:', message.data);
    } else {
        console.log('消息不匹配当前聊天窗口，忽略消息');
        console.log('消息发送者:', sender, '期望:', chattingWith);
        console.log('消息接收者:', receiver, '期望:', myAccount);
        console.log('当前聊天对象:', chattingWith);
        console.log('我的账号:', myAccount);
        
        // 尝试其他可能的字段名
        console.log('尝试其他字段名:');
        console.log('message.senderAccountNum:', message.senderAccountNum);
        console.log('message.receiverAccountNum:', message.receiverAccountNum);
        console.log('message.sender:', message.sender);
        console.log('message.receiver:', message.receiver);
    }
}


// 添加消息到聊天窗口
function AddMessage(message, type = 'received', timestamp = null) {
    const messagesDiv = document.getElementById('messages');

    const p = document.createElement('p');
    p.textContent = message;
    p.classList.add(type);
    messagesDiv.appendChild(p);

    // 添加时间戳
    if (timestamp) {
        const time = document.createElement('div');
        time.textContent = new Date(timestamp).toLocaleTimeString();
        time.classList.add('timestamp');
        p.appendChild(time);
    }

    // 滚动到底部
    scrollToBottom();
}


function updateChatHeader(friend) {
    const chatHeader = document.getElementById("chatHeader");
    if (!friend) {
        chatHeader.innerHTML = `<div class="no-friend-selected">请选择好友开始聊天</div>`;
        return;
    }
    chatHeader.innerHTML = `
        <img src="${friend.avatarUrl || localStorage.getItem('base_url')+'/api/static/images/user-icon.png'}" class="friend-avatar-main" alt="好友头像">
        <div class="friend-info-main">
            <span class="friend-nickname-main">${friend.name}</span>
            <span class="friend-account-main">账号: ${friend.accountNum}</span>
            <span class="friend-status-main">在线</span>
            <span class="friend-desc-main">${friend.desc || "这个人很懒，什么都没写"}</span>
        </div>
        <div class="chat-actions">
            <button class="info-btn" title="好友资料" onclick="showFriendInfo()">i</button>
            <button class="setting-btn" title="聊天设置" onclick="showChatSettings()">⚙</button>
        </div>
    `;
}
function showChatSettings() {
    document.getElementById('infoPanel').style.display = 'block';
    document.getElementById('infoContent').innerHTML = '<h3>聊天设置</h3><p>更多设置...</p>';
}
function goToUserDetailPage() {
    window.location.href = `http://${host}:${port}/api/v1/user/info/start`;
}




// 全局变量存储键盘事件监听器
let kickModalKeyHandler = null;

// 全局变量存储每个好友的对话内容
let friendDraftMessages = {};

// 当前选中的好友（从friend.js中获取）
function getCurrentFriend() {
    return window.currentFriend || null;
}

function setCurrentFriend(friendAccount) {
    window.currentFriend = friendAccount;
}

// 全局字符计数器更新函数
function updateCharCounter(length) {
    const charCounter = document.getElementById('charCounter');
    const MAX_MESSAGE_LENGTH = 1000;
    if (charCounter) {
        charCounter.textContent = `${length}/${MAX_MESSAGE_LENGTH}`;
        
        // 根据长度设置颜色
        charCounter.classList.remove('warning', 'danger');
        if (length > MAX_MESSAGE_LENGTH * 0.9) {
            charCounter.classList.add('danger');
        } else if (length > MAX_MESSAGE_LENGTH * 0.7) {
            charCounter.classList.add('warning');
        }
    }
}

// 滚动到消息区域底部的函数
function scrollToBottom() {
    const messageArea = document.querySelector('.message-area');
    if (messageArea) {
        // 使用平滑滚动效果
        messageArea.scrollTo({
            top: messageArea.scrollHeight,
            behavior: 'smooth'
        });
    }
}

// 将函数设为全局，让其他文件也能访问
window.scrollToBottom = scrollToBottom;

// 获取当前用户的WebSocket连接
function getCurrentWebSocket() {
    const currentUser = getCurrentUserAccount();
    return currentUser ? websocketConnections.get(currentUser) : null;
}

// 将WebSocket获取函数设为全局
window.getCurrentWebSocket = getCurrentWebSocket;

// 清理特定用户的WebSocket连接
function cleanupUserConnection(userAccount) {
    // 关闭WebSocket连接
    const ws = websocketConnections.get(userAccount);
    if (ws) {
        ws.close(1000, "用户切换");
        websocketConnections.delete(userAccount);
    }
    
    // 清除重连定时器
    const timer = reconnectTimers.get(userAccount);
    if (timer) {
        clearTimeout(timer);
        reconnectTimers.delete(userAccount);
    }
    
    // 清除重试计数
    retryCounts.delete(userAccount);
    
    console.log(`已清理用户 ${userAccount} 的WebSocket连接`);
}

// 将清理函数设为全局
window.cleanupUserConnection = cleanupUserConnection;

// 保存当前好友的对话内容
function saveCurrentFriendDraft() {
    if (getCurrentFriend()) {
        const textInput = document.getElementById('textInput');
        if (textInput) {
            friendDraftMessages[getCurrentFriend()] = textInput.value;
            console.log(`保存好友 ${getCurrentFriend()} 的对话内容:`, textInput.value);
            
            // 同时保存到localStorage
            try {
                localStorage.setItem('friendDraftMessages', JSON.stringify(friendDraftMessages));
            } catch (e) {
                console.error('保存对话内容到localStorage失败:', e);
            }
        }
    }
}

// 加载指定好友的对话内容
function loadFriendDraft(friendAccount) {
    const textInput = document.getElementById('textInput');
    if (textInput) {
        const draftMessage = friendDraftMessages[friendAccount] || '';
        textInput.value = draftMessage;
        
        // 更新字符计数器
        const charCounter = document.getElementById('charCounter');
        if (charCounter) {
            charCounter.textContent = `${draftMessage.length}/1000`;
            charCounter.classList.remove('warning', 'danger');
            if (draftMessage.length > 900) {
                charCounter.classList.add('danger');
            } else if (draftMessage.length > 700) {
                charCounter.classList.add('warning');
            }
        }
        
        // 调整输入框高度
        if (typeof adjustHeight === 'function') {
            adjustHeight();
        }
        
        console.log(`加载好友 ${friendAccount} 的对话内容:`, draftMessage);
    }
}

// 清空指定好友的对话内容
function clearFriendDraft(friendAccount) {
    if (friendDraftMessages[friendAccount]) {
        delete friendDraftMessages[friendAccount];
    }
}

function KickUser(message) {
    // 1. 验证当前用户
    const currentUser = sessionStorage.getItem("account_num");
    if (message.accountNum !== currentUser) {
        console.log(`收到踢出消息，但目标不是当前用户（${message.accountNum}）`);
        return;
    }

    // 正确关闭 WebSocket 连接（小写 close，添加特定代码和原因）
    if (ws && ws.readyState !== WebSocket.CLOSED) {
        ws.close(1000, "用户被踢出"); // 1000 表示正常关闭
    }

    // 2. 显示优化的弹窗（保持原有逻辑不变）
    const modalOverlay = document.getElementById("kickModalOverlay");
    const modal = document.getElementById("kickModal");
    const msgElement = document.getElementById("kickMessage");
    const reasonElement = document.getElementById("kickReason");
    const actionBtn = document.getElementById("kickConfirmBtn");

    // 设置消息内容
    msgElement.textContent = "您已在其他设备登录，当前会话已终止。";
    reasonElement.textContent = `原因：${message.data || "新设备登录"}`;

    // 显示弹窗
    modalOverlay.style.display = "block";
    modal.style.display = "block";

    // 添加震动效果
    setTimeout(() => {
        modal.classList.add("shake");
    }, 100);

    // 移除震动效果
    setTimeout(() => {
        modal.classList.remove("shake");
    }, 700);

    // 禁用页面滚动
    document.body.style.overflow = "hidden";

    // 倒计时自动跳转功能
    let countdown = 10;
    const originalBtnText = actionBtn.textContent;

    const updateCountdown = () => {
        if (countdown > 0) {
            actionBtn.textContent = `重新登录 (${countdown}s)`;
            countdown--;
            setTimeout(updateCountdown, 1000);
        } else {
            actionBtn.textContent = originalBtnText;
            confirmKick();
        }
    };

    // 开始倒计时
    setTimeout(updateCountdown, 1000);

    // 添加键盘事件监听（ESC键关闭）
    kickModalKeyHandler = function(event) {
        if (event.key === "Escape") {
            confirmKick();
        }
    };
    document.addEventListener("keydown", kickModalKeyHandler);
}

// 用户点击"重新登录"后执行跳转
function confirmKick() {
    // 移除键盘事件监听器
    if (kickModalKeyHandler) {
        document.removeEventListener("keydown", kickModalKeyHandler);
        kickModalKeyHandler = null;
    }
    
    // 隐藏弹窗
    const modalOverlay = document.getElementById("kickModalOverlay");
    const modal = document.getElementById("kickModal");
    modalOverlay.style.display = "none";
    modal.style.display = "none";
    
    // 恢复页面滚动
    document.body.style.overflow = "";
    
    // 清除登录状态
    localStorage.removeItem("currentUser");
    localStorage.removeItem("token");
    sessionStorage.clear();

    // 关闭 WebSocket 连接
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close(1000, "用户被踢出");
    }

    // 跳转登录页
    window.location.href = `http://${host}:${port}/api/v1/start`;
}

// 防止点击弹窗内容时关闭弹窗
function preventModalClose(event) {
    event.stopPropagation();
}

// 初始化弹窗事件监听器
function initKickModalEvents() {
    const modalOverlay = document.getElementById("kickModalOverlay");
    const modal = document.getElementById("kickModal");
    const confirmBtn = document.getElementById("kickConfirmBtn");
    
    // 点击遮罩层关闭弹窗
    modalOverlay.addEventListener("click", confirmKick);
    
    // 点击弹窗内容不关闭弹窗
    modal.addEventListener("click", preventModalClose);
    
    // 点击确认按钮
    confirmBtn.addEventListener("click", confirmKick);
}

// 页面加载完成后初始化事件
document.addEventListener("DOMContentLoaded", function() {
    updateChatHeader(null);
    initKickModalEvents();
    
    // 初始化对话内容存储
    console.log('初始化对话内容存储功能');
    
    // 从localStorage恢复对话内容（可选，用于页面刷新后保持内容）
    const savedDrafts = localStorage.getItem('friendDraftMessages');
    if (savedDrafts) {
        try {
            friendDraftMessages = JSON.parse(savedDrafts);
            console.log('从localStorage恢复对话内容:', friendDraftMessages);
        } catch (e) {
            console.error('恢复对话内容失败:', e);
            friendDraftMessages = {};
        }
    }
});

// 自动调整输入框高度的功能
function initAutoResizeInput() {
    const textInput = document.getElementById('textInput');
    
    if (!textInput) {
        console.error('找不到textInput元素');
        return;
    }
    
    console.log('开始初始化自动调整输入框高度功能');
    
    // 添加发送标志和消息长度限制
    let isSending = false;
    const MAX_MESSAGE_LENGTH = 1000; // 最大消息长度1000字符
    
    // 创建隐藏的div来计算文本高度
    const hiddenDiv = document.createElement('div');
    hiddenDiv.style.cssText = `
        position: absolute;
        top: -9999px;
        left: -9999px;
        width: ${textInput.offsetWidth}px;
        font-size: 16px;
        font-family: Arial, sans-serif;
        line-height: 1.4;
        padding: 15px;
        border: 1px solid #ccc;
        box-sizing: border-box;
        white-space: pre-wrap;
        word-wrap: break-word;
        overflow: hidden;
        min-height: 20px;
        max-height: none;
    `;
    document.body.appendChild(hiddenDiv);
    
    // 调整输入框高度的函数
    function adjustHeight() {
        const text = textInput.value;
        
        // 更新字符计数器
        const charCounter = document.getElementById('charCounter');
        if (charCounter) {
            charCounter.textContent = `${text.length}/1000`;
            charCounter.classList.remove('warning', 'danger');
            if (text.length > 900) {
                charCounter.classList.add('danger');
            } else if (text.length > 700) {
                charCounter.classList.add('warning');
            }
        }
        
        // 检查消息长度
        if (text.length > MAX_MESSAGE_LENGTH) {
            textInput.value = text.substring(0, MAX_MESSAGE_LENGTH);
            console.log(`消息长度超过限制，已截断到${MAX_MESSAGE_LENGTH}字符`);
            // 更新字符计数器为截断后的长度
            if (charCounter) {
                charCounter.textContent = `${MAX_MESSAGE_LENGTH}/1000`;
                charCounter.classList.add('danger');
            }
            return;
        }
        
        hiddenDiv.textContent = text || 'A'; // 使用'A'作为最小高度参考
        
        // 强制重新计算布局
        hiddenDiv.offsetHeight;
        
        const scrollHeight = hiddenDiv.scrollHeight;
        const newHeight = Math.min(Math.max(scrollHeight, 20), 400);
        textInput.style.setProperty('height', newHeight + 'px', 'important');
        

        
        // 更新隐藏div的宽度（以防窗口大小改变）
        hiddenDiv.style.width = textInput.offsetWidth + 'px';
    }
    
    // 监听输入事件
    textInput.addEventListener('input', function(e) {
        adjustHeight();
        // 保存当前好友的对话内容
        saveCurrentFriendDraft();
    });
    
    // 监听keyup事件作为备用
    textInput.addEventListener('keyup', function(e) {
        adjustHeight();
        // 保存当前好友的对话内容
        saveCurrentFriendDraft();
    });
    
    // 监听keydown事件，处理Enter键
    textInput.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            if (e.ctrlKey) {
                // Ctrl+Enter 发送消息
                e.preventDefault();
                isSending = true;
                const form = document.getElementById('textInputForm');
                if (form) {
                    form.dispatchEvent(new Event('submit'));
                }
                setTimeout(() => { isSending = false; }, 200);
            } else {
                // 普通Enter键，插入换行符
                e.preventDefault();
                const start = this.selectionStart;
                const end = this.selectionEnd;
                const value = this.value;
                
                // 检查插入换行符后是否会超过长度限制
                if (value.length + 1 > MAX_MESSAGE_LENGTH) {
                    console.log('插入换行符会超过长度限制');
                    return;
                }
                
                this.value = value.substring(0, start) + '\n' + value.substring(end);
                this.selectionStart = this.selectionEnd = start + 1;
                adjustHeight();
            }
        }
    });
    
    // 监听聚焦事件
    textInput.addEventListener('focus', function() {
        adjustHeight();
        // 聚焦时滚动到输入框底部
        setTimeout(() => {
            textInput.scrollTop = textInput.scrollHeight;
        }, 100);
    });
    
    // 监听失焦事件
    textInput.addEventListener('blur', function(e) {
        // 延迟处理失焦，避免与发送按钮冲突
        setTimeout(() => {
            // 如果正在发送消息，不收缩输入框
            if (isSending) {
                return;
            }
            
            // 检查是否是因为点击发送按钮而失焦
            const relatedTarget = e.relatedTarget;
            if (relatedTarget && (relatedTarget.classList.contains('send-btn') || relatedTarget.type === 'submit')) {
                return; // 不收缩输入框
            }
            
            // 检查当前焦点是否还在输入框上
            if (document.activeElement === textInput) {
                return;
            }
            
            // 失焦时恢复最小高度
            textInput.style.setProperty('height', '20px', 'important');
        }, 50);
    });
    
    // 监听窗口大小改变事件
    window.addEventListener('resize', adjustHeight);
    
    // 初始化高度
    adjustHeight();
    
    // 初始化字符计数器
    const charCounter = document.getElementById('charCounter');
    if (charCounter) {
        charCounter.textContent = '0/1000';
        charCounter.classList.remove('warning', 'danger');
    }
    

}

// 页面加载完成后初始化自动调整功能
document.addEventListener('DOMContentLoaded', function() {
    initAutoResizeInput();
});

// 如果DOM已经加载完成，直接初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initAutoResizeInput);
} else {
    initAutoResizeInput();
}
