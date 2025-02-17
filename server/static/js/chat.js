

const account_num = new URLSearchParams(window.location.search).get("account_num");

let ws
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
        // 创建 WebSocket 连接
        ws =  new WebSocket(`ws://${host}:${port}/api/v1/chat/ws?token=${localStorage.getItem('jwt')}`);
        // WebSocket 事件监听
        ws.onopen = function (event) {
            console.log('连接已打开', event);
            addMessage('已连接到服务器。', 'sent');
        };
        // 接收消息
        ws.onmessage = function (event) {
            console.log(event.data);
            const message = JSON.parse(event.data);
            switch (message.handlerName) {
                case "RefreshText":
                    break;
                case "SendText":
                    sendTextResponse(message);
                    break;
                case "Login":
                    break;
                case "Logoff":
                    break;
                case "Register":
                    break;
            }
        };
        ws.onclose = function (event) {
            logoff();
            console.log('连接已关闭', event);
            addMessage('已从服务器断开。', 'sent');
            logoff();
        };

        ws.onerror = function (error) {
            console.error('WebSocket 错误', error);
        };
        // 发送消息
        document.getElementById("textInputForm").addEventListener("submit", function (event) {
            event.preventDefault();
            const input = document.getElementById("textInput");
            const message = input.value;

            if (!currentFriend) {
                alert("请先选择一个好友！");
                return;
            }

            const userInput = document.getElementById('textInput').value;
            if (!userInput.trim()) return; // 如果输入为空，则不发送

            const data = JSON.stringify({
                SenderAccountNum: account_num,
                ReceiverAccountNum: currentFriend,
                handlerName: "SendText",
                content: userInput,
            });

            const jsonData = JSON.stringify({
                data: data,
                accountNum: account_num,
                receiver: currentFriend,
                type: "SendText",
            });
            ws.send(jsonData);
            // 显示消息
            addMessage(message,'sent')
            input.value = "";
        });
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });


// 接收到sendText响应
function sendTextResponse(message) {
    if (message.SenderAccountNum === currentFriend) {
        if (message.handlerName === "error") {
            addMessage(`错误: ${message.content}`, 'error');
        }else {
            addMessage(message.content)
        }
    }
}


// 添加消息到聊天窗口
function addMessage(message, type = 'received', timestamp = null) {
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
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}


