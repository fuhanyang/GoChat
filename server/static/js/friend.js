let currentFriend = null; // 当前选中的好友
// 假设好友列表从服务器获取
let friends = []; // 示例好友列表
fetch('../static/config.json')
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(config => {
        // 初始化应用程序，传入配置信息
        initializeApp(config);
        // 初始化
        loadFriends();
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
function displayFriendList() {
    const friendList = document.getElementById("friendList");
    friendList.innerHTML = ""; // 清空列表
    friends.forEach(friend => {
        const li = document.createElement("li");
        li.textContent = friend.name;
        const span = document.createElement("span");
        span.textContent = friend.accountNum;
        li.appendChild(span);
        li.addEventListener("click", () => switchFriend(friend));
        friendList.appendChild(li);
    });
}
// 加载好友列表
function loadFriends() {
    fetch(`http://${host}:${port}/api/v1/friend/list`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify({
            accountNum  : account_num,
            handlerName :"GetFriends",
        })
    }).then(response => response.json()
    ).then(data => {
        if(data.error){
            throw data.error;
        }
        data.data.forEach(friend => {
            friends.push(friend);
            console.log(friend.accountNum+friend.name+"load");
            displayFriendList();
        })
    }).catch(error => {
        console.log(error);
        alert("load friends error"+error);
    })
} // 切换好友
function switchFriend(friend) {
    console.log("切换到好友", friend.accountNum, friend.name);
    currentFriend = friend.accountNum;
    loadMessages();
    // 更新用户信息显示区域
    const userInfo = document.getElementById("userInfo");
    userInfo.innerHTML = `
        <h4>当前好友信息</h4>
        <p>账号: ${friend.accountNum}</p>
        <p>昵称: ${friend.name}</p>
    `;
}
function addFriend() {
    fetch(`http://${host}:${port}/api/v1/friend/addition`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify({
            accountNum: account_num,
            handlerName: "AddFriend",
        })
    }).then(response => {
        return response.json();
    }).then(data => {
        if(data.error){
            throw data.error;
        }
        console.log("添加好友成功",data.data.accountNum,data.data.name);
        let friend = {
            accountNum: data.data.accountNum,
            name: data.data.name,
        }
        friends.push(friend);
       displayFriendList();
    }).catch(error => {
        alert("添加好友失败，请重试。"+error);
    });
}
// 加载消息记录
function loadMessages() {
    const messagesDiv = document.getElementById("messages");
    messagesDiv.innerHTML = ""; // 清空消息
    addMessage(`正在加载与${currentFriend}的对话记录...`, 'loading');

    // 构造请求体
    const requestBody = {
        SenderAccountNum: account_num,
        ReceiverAccountNum: currentFriend,
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
        .then(response => {
            if (!response.ok) {
                // 如果状态码不是 2xx，抛出错误
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            messagesDiv.innerHTML = ""; // 清空消息
            // 解析并显示消息
            if(data.text.content===null){
                console.log("对话记录为空")
                return
            }
            data.text.content.reverse().forEach(content => {
                let p; // 声明一个变量来存储将要创建的段落元素
                if (content.senderAccountNum === account_num && content.receiverAccountNum === currentFriend) {
                    // 如果消息是由当前用户发送给当前朋友的
                    p = document.createElement("p");
                    p.textContent = `${content.content}`; // 假设消息内容在 content 字段中，并添加前缀“你:”
                    p.classList.add("sent");
                } else if (content.senderAccountNum === currentFriend && content.receiverAccountNum === account_num) {
                    // 如果消息是由当前朋友发送给当前用户的
                    p = document.createElement("p");
                    p.textContent = `${content.content}`; // 假设有 senderName 字段，或默认使用“对方”作为前缀
                    p.classList.add("received");
                } else {
                    // 如果消息的发送者或接收者不是预期的人，可以选择不创建段落元素或打印错误
                    console.error("消息内容格式错误或不属于当前对话");
                    return; // 退出当前迭代，不添加任何元素到 DOM 中
                }
                // 只有当满足条件时才将段落元素添加到 DOM 中
                messagesDiv.appendChild(p);
            });
        })
        .catch(error => {
            console.error('Error:', error);
            addMessage('加载消息失败，请重试。', 'error');
        });
}

