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

window.currentFriend = null; // 当前选中的好友，设为全局变量
// 假设好友列表从服务器获取
let friends = []; // 示例好友列表

loadFriends();
// function displayFriendList() {
//     Const friendList = document.getElementById("friendList");
//     friendList.innerHTML = ""; // 清空列表
//     friends.forEach(friend => {
//         Const li = document.createElement("li");
//         li.textContent = friend.name;
//         Const span = document.createElement("span");
//         span.textContent = friend.accountNum;
//         li.appendChild(span);
//         li.addEventListener("click", () => switchFriend(friend));
//         friendList.appendChild(li);
//     });
// }
function displayFriendList() {
    const friendList = document.getElementById("friendList");
    friendList.innerHTML = "";
    
    if (friends.length === 0) {
        // 显示空状态
        const emptyDiv = document.createElement("div");
        emptyDiv.className = "friend-list-empty";
        emptyDiv.innerHTML = `
            <div class="friend-list-empty-icon">👥</div>
            <div class="friend-list-empty-text">暂无好友</div>
            <div class="friend-list-empty-hint">点击上方按钮添加好友开始聊天吧</div>
        `;
        friendList.appendChild(emptyDiv);
        return;
    }
    
    friends.forEach(friend => {
        const li = document.createElement("li");
        li.className = (window.currentFriend === friend.accountNum) ? "active" : "";
        // 头像
        const img = document.createElement("img");
        img.src = friend.avatarUrl || localStorage.getItem('base_url')+"/api/static/images/user-icon.png";
        img.className = "friend-avatar";
        // 信息
        const infoDiv = document.createElement("div");
        infoDiv.className = "friend-info-list";
        infoDiv.innerHTML = `
                <span class="friend-nickname-list">${friend.name}</span>
                <span class="friend-account-list">账号: ${friend.accountNum}</span>
                <span class="friend-desc-list">${friend.desc || "这个人很懒，什么都没写"}</span>
            `;
        li.appendChild(img);
        li.appendChild(infoDiv);
        li.addEventListener("click", () => switchFriend(friend));
        friendList.appendChild(li);
    });
}
// 加载好友列表
function loadFriends() {
    fetch(localStorage.getItem('base_url')+`/api/v1/friend/list`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify({
            accountNum  : getCurrentUserAccount(),
            handlerName :"GetFriends",
        })
    }).then(response => response.json()
    ).then(data => {
        if(data.error){
            throw data.error;
        }
        if(data.code!==200){
            showNotification('error', '加载失败', data.msg, '操作结果');
            throw new Error(data.msg);
        }
        // 清空现有好友列表
        friends = [];
        // 添加所有好友
        data.data.forEach(friend => {
            friends.push(friend);
            console.log(friend.accountNum+friend.name+"load");
        });
        // 只调用一次displayFriendList
        displayFriendList();
    }).catch(error => {
        console.log(error);
        showNotification('error', '加载失败', "load friends error: " + error, '网络错误');
    })
} // 切换好友
function switchFriend(friend) {
    console.log("切换到好友", friend.accountNum, friend.name);
    
    // 保存当前好友的对话内容
    if (window.currentFriend && typeof saveCurrentFriendDraft === 'function') {
        saveCurrentFriendDraft();
    }
    
    // 切换到新好友
    window.currentFriend = friend.accountNum;
    
    // 加载新好友的对话内容
    if (typeof loadFriendDraft === 'function') {
        loadFriendDraft(friend.accountNum);
    }
    
    loadMessages();
    displayFriendList();
    
    // 切换好友后滚动到底部
            if (typeof scrollToBottom === 'function') {
            setTimeout(() => {
                scrollToBottom();
            }, 200);
        }
}
// 显示友好提示弹窗
function showNotification(type, title, message, subtitle = "系统消息") {
    const modal = document.getElementById('notificationModal');
    const overlay = document.getElementById('notificationModalOverlay');
    const icon = document.getElementById('notificationIcon');
    const titleEl = document.getElementById('notificationTitle');
    const subtitleEl = document.getElementById('notificationSubtitle');
    const messageEl = document.getElementById('notificationMessage');
    
    // 设置弹窗类型和样式
    modal.className = `notification-modal ${type}`;
    
    // 根据类型设置图标
    switch(type) {
        case 'warning':
            icon.textContent = '⚠️';
            break;
        case 'info':
            icon.textContent = 'ℹ️';
            break;
        case 'success':
            icon.textContent = '✅';
            break;
        case 'error':
            icon.textContent = '❌';
            break;
        default:
            icon.textContent = 'ℹ️';
    }
    
    // 设置内容
    titleEl.textContent = title;
    subtitleEl.textContent = subtitle;
    messageEl.textContent = message;
    
    // 显示弹窗
    overlay.style.display = 'block';
    modal.style.display = 'block';
    
    // 添加确认按钮事件
    const confirmBtn = document.getElementById('notificationConfirmBtn');
    confirmBtn.onclick = function() {
        hideNotification();
    };
    
    // 点击遮罩层关闭弹窗
    overlay.onclick = function() {
        hideNotification();
    };
    
    // 3秒后自动关闭（可选）
    setTimeout(() => {
        if (overlay.style.display === 'block') {
            hideNotification();
        }
    }, 3000);
}

// 隐藏友好提示弹窗
function hideNotification() {
    const modal = document.getElementById('notificationModal');
    const overlay = document.getElementById('notificationModalOverlay');
    overlay.style.display = 'none';
    modal.style.display = 'none';
}

function addFriend() {
    // 获取搜索框中的内容
    const searchInput = document.getElementById('friendSearch');
    const searchValue = searchInput.value.trim();
    
    // 如果搜索框有内容，先检查本地好友列表
    if (searchValue) {
        // 检查是否已经是好友
        const existingFriend = friends.find(friend => friend.accountNum === searchValue);
        if (existingFriend) {
            showNotification('warning', '好友已存在', `该账号 ${searchValue} 已经是您的好友了！`, '重复添加');
            return;
        }
        
        // 检查是否是自己
        if (searchValue === getCurrentUserAccount()) {
            showNotification('error', '操作无效', '不能添加自己为好友！', '无效操作');
            return;
        }
    }
    
    // 根据搜索框内容决定调用哪个API
    let requestBody = {
        accountNum: getCurrentUserAccount(),
    };
    
    let apiUrl;
    
    if (searchValue) {
        // 如果搜索框有内容，调用AddFriendWithAccountNum
        apiUrl = `http://${host}:${port}/api/v1/friend/addition/with_account_num`;
        requestBody.handlerName = "AddFriendWithAccountNum";
        requestBody.targetAccountNum = searchValue;
        console.log("通过账号添加好友:", searchValue);
    } else {
        // 如果搜索框为空，调用AddFriend随机匹配
        apiUrl = `http://${host}:${port}/api/v1/friend/addition`;
        requestBody.handlerName = "AddFriend";
        console.log("随机匹配好友");
    }
    
    fetch(apiUrl, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(requestBody)
    }).then(async response => {
        const data = await response.json();
        if(!response.ok){
            showNotification('error', '添加失败', data.msg, '操作结果');
            return
        }
        console.log("添加好友成功",data.data.accountNum,data.data.name);
        let friend = {
            accountNum: data.data.accountNum,
            name: data.data.name,
        }
        friends.push(friend);
        displayFriendList();
        
        // 清空搜索框
        searchInput.value = '';
        
        // 显示成功提示
        showNotification('success', '添加成功', `已成功添加 ${friend.name} (${friend.accountNum}) 为好友！`, '操作结果');
    }).catch(error => {
        showNotification('error', '添加失败', `添加好友失败，请重试。\n${error}`, '网络错误');
    });
}
// 加载消息记录
function loadMessages() {
    const messagesDiv = document.getElementById("messages");
    messagesDiv.innerHTML = ""; // 清空消息
    AddMessage(`正在加载与${window.currentFriend}的对话记录...`, 'loading');

    // 构造请求体
    const requestBody = {
        SenderAccountNum: getCurrentUserAccount(),
        ReceiverAccountNum: window.currentFriend,
        handlerName  :"RefreshText",
    };

    // 发送 POST 请求到服务器
    fetch(`http://${host}:${port}/api/v1/chat/text/refresh`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(
            requestBody
        )
    })
        .then(async response => {
            const data = await response.json();
            if (!response.ok) {
                alert(data.msg);
                return
            }
            messagesDiv.innerHTML = ""; // 清空消息
            // 解析并显示消息
            if(data.text.content===null){
                console.log("对话记录为空")
                // 即使没有消息，也要滚动到底部
                if (typeof scrollToBottom === 'function') {
                    scrollToBottom();
                }
                return
            }
            data.text.content.reverse().forEach(content => {
                let p; // 声明一个变量来存储将要创建的段落元素
                if (content.senderAccountNum === getCurrentUserAccount() && content.receiverAccountNum === window.currentFriend) {
                    // 如果消息是由当前用户发送给当前朋友的
                    p = document.createElement("p");
                    p.textContent = `${content.content}`; // 假设消息内容在 content 字段中，并添加前缀"你:"
                    p.classList.add("sent");
                } else if (content.senderAccountNum === window.currentFriend && content.receiverAccountNum === getCurrentUserAccount()) {
                    // 如果消息是由当前朋友发送给当前用户的
                    p = document.createElement("p");
                    p.textContent = `${content.content}`; // 假设有 senderName 字段，或默认使用"对方"作为前缀
                    p.classList.add("received");
                } else {
                    // 如果消息的发送者或接收者不是预期的人，可以选择不创建段落元素或打印错误
                    console.error("消息内容格式错误或不属于当前对话");
                    return; // 退出当前迭代，不添加任何元素到 DOM 中
                }
                // 只有当满足条件时才将段落元素添加到 DOM 中
                messagesDiv.appendChild(p);
            });
            
            // 加载完消息后自动滚动到底部
            if (typeof scrollToBottom === 'function') {
                // 使用setTimeout确保DOM更新完成后再滚动
                setTimeout(() => {
                    scrollToBottom();
                }, 100);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert(error)
            AddMessage('加载消息失败，请重试。', 'error');
            // 即使加载失败，也要尝试滚动到底部
            if (typeof scrollToBottom === 'function') {
                scrollToBottom();
            }
        });
}

document.addEventListener("DOMContentLoaded", function() {
    updateChatHeader(null);
    
    // 初始化添加好友按钮的默认文本
    const addFriendBtn = document.querySelector('.add-friend-btn');
    if (addFriendBtn) {
        addFriendBtn.textContent = '随机匹配好友';
    }
    

});

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
        `;
}
function showFriendInfo() {
    document.getElementById('infoPanel').style.display = 'block';
    document.getElementById('infoContent').innerHTML = '<h3>好友资料</h3><p>更多信息...</p>';
}

