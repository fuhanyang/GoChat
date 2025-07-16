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
// 读取
const savedAccountNum = getCurrentUserAccount();
document.getElementById('userAccountNum').textContent = '账号：' + (savedAccountNum || '');
// 获取用户信息
const jwt = getCurrentUserJWT();

fetch(localStorage.getItem('base_url')+`/api/v1/user/info?account_num=${savedAccountNum}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': jwt ? `Bearer ${jwt}` : ''
        },
    })
    .then(res => res.json())
    .then(data => {
        if (data.code !== 200) {
            document.getElementById('userNickname').textContent = '用户信息加载失败';
            document.getElementById('userDesc').textContent = data.msg || '';
            return;
        }
        const user = data.data;
        document.getElementById('userNickname').textContent = user.name || '未设置昵称';
        document.getElementById('userAccountNum').textContent = '账号：' + (user.account_num || "");
        document.getElementById('userDesc').textContent = user.desc || '这个人很懒，什么都没写';
        document.getElementById('userAvatar').src = getUserAvatar();
        if (user.create_at) {
            document.getElementById('userRegisterTime').textContent = '注册时间：' + new Date(user.create_at).toLocaleString();
        }
        // 更多信息
        document.getElementById('userEmail').textContent = '邮箱：' + (user.email || '-');
        document.getElementById('userPhone').textContent = '手机号：' + (user.phone || '-');
        document.getElementById('userLastLogin').textContent = '最近登录：' + (user.lastLogin || '-');
        document.getElementById("ip").textContent = user.ip || '-';
    })
    .catch(err => {
        document.getElementById('userNickname').textContent = '用户信息加载失败';
        document.getElementById('userDesc').textContent = err.message;
    });

// 头像上传功能
function initAvatarUpload() {
    // 创建隐藏的文件输入框
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = 'image/*';
    fileInput.style.display = 'none';
    document.body.appendChild(fileInput);

    // 为头像上传按钮添加点击事件
    const uploadBtn = document.getElementById('avatarUploadBtn');
    if (uploadBtn) {
        uploadBtn.addEventListener('click', () => {
            fileInput.click();
        });
    }

    // 文件选择处理
    fileInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (!file) return;

        // 验证文件类型和大小
        if (!file.type.startsWith('image/')) {
            showMessage('请选择图片文件', 'error');
            return;
        }

        if (file.size > 5 * 1024 * 1024) { // 5MB限制
            showMessage('图片大小不能超过5MB', 'error');
            return;
        }

        // 预览图片
        previewImage(file);
    });
}

// 预览图片
function previewImage(file) {
    const reader = new FileReader();
    reader.onload = (e) => {
        const userAvatar = document.getElementById('userAvatar');
        userAvatar.src = e.target.result;

        // 显示上传确认对话框
        showUploadConfirm(file);
    };
    reader.readAsDataURL(file);
}

// 显示上传确认对话框
function showUploadConfirm(file) {
    // 创建确认对话框
    const modal = document.createElement('div');
    modal.className = 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50';
    modal.innerHTML = `
        <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 class="text-lg font-semibold mb-4">确认上传头像</h3>
            <div class="flex items-center space-x-4 mb-6">
                <img src="${URL.createObjectURL(file)}" alt="预览" class="w-16 h-16 rounded-full object-cover">
                <div>
                    <p class="text-sm text-gray-600">文件名: ${file.name}</p>
                    <p class="text-sm text-gray-600">大小: ${(file.size / 1024 / 1024).toFixed(2)}MB</p>
                </div>
            </div>
            <div class="flex justify-end space-x-3">
                <button onclick="this.closest('.fixed').remove()" class="px-4 py-2 text-gray-600 hover:text-gray-800">取消</button>
                <button onclick="uploadAvatar(this, '${file.name}')" class="px-4 py-2 bg-primary text-white rounded hover:bg-blue-600">确认上传</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
}

// 上传头像到服务器
async function uploadAvatar(button, fileName) {
    const fileInput = document.querySelector('input[type="file"]');
    const file = fileInput.files[0];
    if (!file) return;

    // 显示上传中状态
    button.disabled = true;
    button.textContent = '上传中...';
    button.classList.add('opacity-50');

    try {
        // 获取配置
        const configResponse = await fetch(localStorage.getItem('base_url')+'/api/static/config.json');
        const config = await configResponse.json();

        // 创建FormData
        const formData = new FormData();
        formData.append('avatar', file);
        formData.append('account_num', savedAccountNum);

        // 获取JWT
        const jwt = getCurrentUserJWT();

        // 上传请求
        const response = await fetch(`http://${config.host}:${config.port}/api/v1/user/upload-avatar`, {
            method: 'POST',
            headers: {
                'Authorization': jwt ? `Bearer ${jwt}` : ''
            },
            body: formData
        });

        const result = await response.json();

        if (result.code === 200) {
            showMessage('头像上传成功！', 'success');
            // 更新头像URL
            document.getElementById('userAvatar').src = result.data.avatarUrl || result.data.url;
            localStorage.setItem('avatar_url', result.data.avatarUrl || result.data.url)
            // 关闭对话框
            document.querySelector('.fixed').remove();
        } else {
            showMessage(result.msg || '上传失败', 'error');
        }
    } catch (error) {
        console.error('上传错误:', error);
        showMessage('上传失败，请重试', 'error');
    } finally {
        // 恢复按钮状态
        button.disabled = false;
        button.textContent = '确认上传';
        button.classList.remove('opacity-50');
    }
}

// 显示消息提示
function showMessage(message, type = 'info') {
    // 创建消息提示元素
    const messageDiv = document.createElement('div');
    messageDiv.className = `fixed top-4 right-4 z-50 px-6 py-3 rounded-lg shadow-lg transition-all duration-300 transform translate-x-full`;

    // 根据类型设置样式
    if (type === 'success') {
        messageDiv.className += ' bg-green-500 text-white';
    } else if (type === 'error') {
        messageDiv.className += ' bg-red-500 text-white';
    } else {
        messageDiv.className += ' bg-blue-500 text-white';
    }

    messageDiv.textContent = message;
    document.body.appendChild(messageDiv);

    // 显示动画
    setTimeout(() => {
        messageDiv.classList.remove('translate-x-full');
    }, 100);

    // 自动隐藏
    setTimeout(() => {
        messageDiv.classList.add('translate-x-full');
        setTimeout(() => {
            document.body.removeChild(messageDiv);
        }, 300);
    }, 3000);
}

// 页面加载完成后初始化头像上传功能
document.addEventListener('DOMContentLoaded', () => {
    initAvatarUpload();
});

function goBack() {
    window.history.length > 1 ? window.history.back() : window.location.href = 'chat.html';
}

function editUser() {
    const accountNum = new URLSearchParams(window.location.search).get('account_num');
    window.location.href = `editUser.html?account_num=${accountNum}`;
}

// 定义获取头像的函数
function getUserAvatar() {
    if (localStorage.getItem('avatar_url') !== null) {
        return localStorage.getItem('avatar_url');
    } else {
        return localStorage.getItem('base_url') + '/api/static/images/user-icon.png'; // 默认头像
    }
}



