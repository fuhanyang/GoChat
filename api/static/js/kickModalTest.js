// 踢出用户弹窗测试功能
// 仅在开发环境中使用

// 测试弹窗功能
function testKickModal() {
    const testMessage = {
        accountNum: sessionStorage.getItem("account_num") || "test_user",
        data: "测试踢出原因 - 模拟其他设备登录"
    };
    KickUser(testMessage);
}

// 在开发环境中添加测试按钮
function addTestButton() {
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        const testBtn = document.createElement('button');
        testBtn.textContent = '测试踢出弹窗';
        testBtn.style.cssText = `
            position: fixed; 
            bottom: 20px; 
            right: 20px; 
            z-index: 9998; 
            padding: 10px 15px; 
            background: #007bff; 
            color: white; 
            border: none; 
            border-radius: 5px; 
            cursor: pointer;
            font-size: 12px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.2);
        `;
        testBtn.onclick = testKickModal;
        document.body.appendChild(testBtn);
        
        // 添加悬停效果
        testBtn.addEventListener('mouseenter', () => {
            testBtn.style.background = '#0056b3';
            testBtn.style.transform = 'translateY(-1px)';
        });
        
        testBtn.addEventListener('mouseleave', () => {
            testBtn.style.background = '#007bff';
            testBtn.style.transform = 'translateY(0)';
        });
    }
}

// 页面加载完成后添加测试按钮
document.addEventListener("DOMContentLoaded", function() {
    addTestButton();
}); 