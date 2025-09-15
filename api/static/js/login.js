let Password = "";
let host = "localhost"//"8.137.19.65"
let port = "8080"
let base_url
let account_num
(function initializeApp() {
    localStorage.setItem('base_url', `http://${host}:${port}`);
    console.log(`Application initialized with host: ${host} and port: ${port}`);
    account_num = new URLSearchParams(window.location.search).get("account_num");
// 保存到当前页面的 sessionStorage
    if (account_num) {
        sessionStorage.setItem("account_num", account_num);
        console.log("当前页面已保存 account_num：", account_num);
    }
    // 在这里你可以继续初始化应用程序的其他部分
})()

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
            ip:              host,
            handlerName:     "Register",
        }
        ),
    })
    // 注册成功处理
    .then(async response => {
        const data = await response.json();
        if (response.status === 200) {
            document.getElementById('successMessage').textContent = data.data || '注册成功！';
            document.getElementById('successModal').classList.remove('hidden');
            showLogin();
        } else {
            document.getElementById('errorMessage').textContent = data.msg || '未知错误';
            document.getElementById('errorModal').classList.remove('hidden');
        }
    })
    .catch(error => {
        console.error("注册失败:", error);
        document.getElementById('errorMessage').textContent = '注册失败，请重试！' + (error.message || '');
        document.getElementById('errorModal').classList.remove('hidden');
    });
    // .then(async response => {
    //     Const data = await response.json();
    //     if (response.status === 200) {
    //         alert("注册成功！" + data.data);
    //         showLogin();
    //     } else {
    //         alert(data.msg);
    //     }
    // })
    // .catch(error => {
    //     console.error("注册失败:", error);
    //     alert("注册失败，请重试！"+error);
    // });
}

// 登录功能
async function login() {
    try{
        const account_num = document.getElementById("login-account_num").value;
        Password = document.getElementById("login-password").value;
        sessionStorage.setItem("password", Password);
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
                console.log("登录成功！jwt:", jwt);
                // 登录成功后保存 - 使用多用户存储模式
                const userTokens = JSON.parse(localStorage.getItem('userTokens') || '{}');
                userTokens[account_num] = {
                    jwt: jwt,
                    account_num: account_num,
                    loginTime: new Date().toISOString()
                };
                localStorage.setItem('userTokens', JSON.stringify(userTokens));
                sessionStorage.setItem("jwt", jwt);
                // 设置当前活跃用户
                sessionStorage.setItem("currentUser", account_num);
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
            accountNum:  sessionStorage.getItem("account_num"),
            password:    sessionStorage.getItem("password"),
            ip:         "127.0.0.1",
            handlerName: "LogOff",
        }),
    })
        .then(response => {
            if (response.status === 200) {
                window.location.href = `http://${host}:${port}/api/v1/start`;
                localStorage.removeItem('jwt');
            }
        })
        .catch(error => {
            console.error("退出失败:", error);
            alert("退出失败:"+ error)
        });
}


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
