let Password = "";
let host
let port
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
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
function initializeApp(config) {
    // 使用配置信息初始化应用程序
    port = config.port;
    host = config.host;

    console.log(`Application initialized with host: ${host} and port: ${port}`);
    // 在这里你可以继续初始化应用程序的其他部分
}
// 注册功能
function register() {
    const username = document.getElementById("register-username").value;
    const password = document.getElementById("register-password").value;
    const passwordConfirm = document.getElementById("register-passwordConfirm").value;
    fetch(`http://${host}:${port}/api/v1/user/register`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(
            {
            username:        username,
            password:        password,
            passwordConfirm: passwordConfirm,
            ip:              "127.0.0.1",
            handlerName:     "Register",
        }
        ),
    })
        .then(response => {
            if (response.status === 201) {
                alert("注册成功！");
                showLogin();
            } else {
                response.json().then(message => alert(message.data));
            }
        })
        .catch(error => {
            console.error("注册失败:", error);
            alert("注册失败，请重试！");
        });
}

// 登录功能
async function login() {
    try{
        const account_num = document.getElementById("login-account_num").value;
        Password = document.getElementById("login-password").value;
        const response = await fetch(`http://${host}:${port}/api/v1/user/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                accountNum: account_num,
                password: Password,
                ip: `${host}:${port}`,
                handlerName: "Login",
            }),
        })
        const data =  await response.json();
        const jwt = response.headers.get("Authorization");
        if (!response.ok) {
            console.error("登录失败:", data.error);
            alert("登录失败，请重试！" + data.error);
        } else {
            if (jwt) {
                // 存储JWT到localStorage或sessionStorage
                localStorage.setItem('jwt', jwt);
                console.log("登录成功！jwt:", jwt);
                window.location.href = `http://${host}:${port}/api/v1/chat?account_num=${account_num}`;
            } else {
                console.error("登录失败，请重试！")
                alert("jwt获取失败！");
            }
        }
    }catch(error){
        console.error("登录失败:", error);
        alert("登录失败，请重试！" + error.message);
    }

}

// 退出功能
 function logoff() {
    fetch(`http://${host}:${port}/api/v1/user/logoff`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify({
            accountNum:  account_num,
            password:    "123",
            ip:         "127.0.0.1",
            handlerName: "Logoff",
        }),
    })
        .then(response => {
            if (response.status === 200) {
                window.location.href = `http://${host}:${port}/api/v1/start`;
                document.getElementById("logoff").style.display = "none";
                document.getElementById("login-form").style.display = "block";
                document.getElementById("login-username").value = "";
                document.getElementById("login-password").value = "";
            }
        })
        .catch(error => {
            console.error("退出失败:", error);
        });
}
// 模拟用户数据存储（实际项目中应使用后端数据库）
let users = JSON.parse(localStorage.getItem("users")) || [];

// 显示注册表单
function showRegister() {
    document.getElementById("login-form").style.display = "none";
    document.getElementById("register-form").style.display = "block";
}

// 显示登录表单
function showLogin() {
    document.getElementById("register-form").style.display = "none";
    document.getElementById("login-form").style.display = "block";

}