<template>
    <nav class="navbar">
        <!-- 品牌标志区域 -->
        <div class="navbar-brand" @click="$router.push('/')">
            <img src="@/assets/acm-logo.jpg" alt="ACM集训队Logo" class="brand-logo" />
            <div class="brand-text">
                <div class="school-name">成都锦城学院</div>
                <div class="team-name">ACM集训队</div>
            </div>
            <img src="@/assets/acm-icon.png" alt="ACM图标" class="acm-icon" />
        </div>

        <!-- 移动端汉堡菜单按钮 -->
        <div class="hamburger" @click="toggleMobileMenu" :class="{ 'active': mobileMenuActive }">
            <div class="bar"></div>
            <div class="bar"></div>
            <div class="bar"></div>
        </div>

        <!-- 自适应间隔 -->
        <div class="spacer1"></div>

        <!-- 导航链接 -->
        <ul class="nav-links" :class="{ 'mobile-active': mobileMenuActive }">
            <li><router-link to="/" @click="closeMobileMenu">首页</router-link></li>
            <li><router-link to="/knowledge" @click="closeMobileMenu">知识库</router-link></li>
            <li><router-link to="/competition" @click="closeMobileMenu">比赛页面</router-link></li>
            <li><router-link to="/ranking" @click="closeMobileMenu">排行榜</router-link></li>
            <li><router-link to="/faq" @click="closeMobileMenu">常见问题</router-link></li>
        </ul>

        <!-- 自适应间隔 -->
        <div class="spacer2"></div>

        <!-- 登录 / 控制面板 + 退出 -->
        <div v-if="isLoggedIn" class="auth-buttons" :class="{ 'mobile-active': mobileMenuActive }">
            <router-link v-if="role === 'admin'" to="/admin/dashboard" class="control-panel-btn"
                @click="closeMobileMenu">
                管理员控制面板
            </router-link>
            <router-link v-else-if="role === 'member'" to="/member/dashboard" class="control-panel-btn"
                @click="closeMobileMenu">
                队员控制面板
            </router-link>
            <router-link v-else to="/student/dashboard" class="control-panel-btn" @click="closeMobileMenu">
                学生控制面板
            </router-link>
            <button class="logout-btn" @click="handleLogout">退出</button>
        </div>

        <router-link v-else to="/login" class="login-btn" :class="{ 'mobile-active': mobileMenuActive }"
            @click="closeMobileMenu">登录</router-link>
    </nav>
</template>

<script>
import emitter from '@/utils/eventBus'

export default {
    name: 'NavBar',
    data() {
        return {
            role: null,
            isLoggedIn: false,
            mobileMenuActive: false
        }
    },
    mounted() {
        this.syncLoginState()
        // 监听登录状态变化
        emitter.on('loginChange', this.syncLoginState)
        // 监听窗口大小变化
        window.addEventListener('resize', this.handleResize)
        this.handleResize() // 初始化时检查一次
    },
    beforeUnmount() {
        emitter.off('loginChange', this.syncLoginState)
        window.removeEventListener('resize', this.handleResize)
    },
    methods: {
        syncLoginState() {
            const memberToken = localStorage.getItem('member_token')
            const adminToken = localStorage.getItem('admin_token')
            const userToken = localStorage.getItem('user_token')

            if (memberToken) {
                this.role = 'member'
                this.isLoggedIn = true
            } else if (adminToken) {
                this.role = 'admin'
                this.isLoggedIn = true
            } else if (userToken) {
                const userRole = localStorage.getItem('user_role') || 'user'
                this.role = userRole
                this.isLoggedIn = true
            } else {
                this.role = null
                this.isLoggedIn = false
            }
        },
        handleLogout() {
            localStorage.removeItem('member_token')
            localStorage.removeItem('admin_token')
            localStorage.removeItem('user_token')
            localStorage.removeItem('user_role')
            localStorage.removeItem('member_info') // 删除队员信息
            this.syncLoginState()
            emitter.emit('loginChange', { role: null })
            this.$router.push('/login')
        },
        toggleMobileMenu() {
            this.mobileMenuActive = !this.mobileMenuActive
        },
        closeMobileMenu() {
            this.mobileMenuActive = false
        },
        handleResize() {
            // 在窗口大小变化时，如果宽度大于768px且菜单是打开的，则关闭菜单
            if (window.innerWidth > 768 && this.mobileMenuActive) {
                this.mobileMenuActive = false
            }
        }
    }
}
</script>

<style scoped>
/* 基础导航栏样式 */
.navbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 80px;
    background: linear-gradient(90deg, #f8fafc, #e2e8f0, #cbd5e1);
    color: #1e293b;
    display: flex;
    align-items: center;
    padding: 0 2rem;
    z-index: 1000;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.04),
        0 2px 4px -1px rgba(0, 0, 0, 0.02);
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

/* 自适应间隔 */
.spacer1 {
    flex: 3;
}

.spacer2 {
    flex: 4;
}

/* 导航链接样式 */
.nav-links {
    list-style: none;
    display: flex;
    gap: 1.5rem;
    margin: 0;
    padding: 0;
}

.nav-links li {
    position: relative;
    padding: 0.5rem;
}

.nav-links a {
    color: #1e293b;
    text-decoration: none;
    font-weight: 500;
    font-size: 1rem;
    padding: 0.5rem 0;
    transition: all 0.3s ease;
    position: relative;
}

.nav-links a:not(.login-btn):hover {
    color: #3498db;
}

.nav-links a:not(.login-btn)::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 0;
    height: 2px;
    background-color: #3498db;
    transition: width 0.3s ease;
}

.nav-links a:not(.login-btn):hover::after {
    width: 100%;
}

/* 按钮基础样式 */
.login-btn,
.control-panel-btn {
    position: relative;
    background: linear-gradient(135deg, #3B82F6, #3498db);
    color: white !important;
    padding: 0.75rem 1.75rem !important;
    border-radius: 12px;
    font-weight: 600 !important;
    font-size: 1rem;
    text-decoration: none;
    border: none;
    cursor: pointer;
    overflow: hidden;
    box-shadow:
        0 4px 6px rgba(59, 130, 246, 0.25),
        0 1px 3px rgba(59, 130, 246, 0.1),
        inset 0 1px 0 rgba(255, 255, 255, 0.2);
    transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.control-panel-btn::before,
.login-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #3498db, #1D4ED8);
    opacity: 0;
    transition: opacity 0.4s ease;
    z-index: -1;
}

.control-panel-btn:hover,
.login-btn:hover {
    transform: translateY(-3px) scale(1.02);
    box-shadow:
        0 10px 20px rgba(59, 130, 246, 0.3),
        0 6px 6px rgba(59, 130, 246, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.control-panel-btn:hover::before,
.login-btn:hover::before {
    opacity: 1;
}

.control-panel-btn:active,
.login-btn:active {
    transform: translateY(1px) scale(0.98);
    box-shadow:
        0 2px 5px rgba(59, 130, 246, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

/* 品牌标志区域 */
.navbar-brand {
    display: flex;
    align-items: center;
    cursor: pointer;
    transition: all 0.3s ease;
    border-radius: 10px;
}

.navbar-brand:hover {
    background: rgba(0, 0, 0, 0.1);
}

.brand-logo {
    height: 50px;
    width: 50px;
    margin-right: 0.75rem;
    object-fit: cover;
    border-radius: 50%;
    filter: drop-shadow(0 0 6px rgba(100, 116, 139, 0.2));
    transition: transform 0.3s ease;
}

.navbar-brand:hover .brand-logo {
    transform: scale(1.05);
}

.brand-text {
    line-height: 1.2;
}

.school-name {
    font-size: 14px;
    color: #475569;
    font-weight: 500;
    display: flex;
    margin-left: 5px;
    justify-content: flex-start;
    text-align: left;
}

.team-name {
    font-size: 1.0625rem;
    font-weight: 700;
    color: #1e293b;
    letter-spacing: 0.025rem;
    text-shadow: 0 0 6px rgba(100, 116, 139, 0.2);
    transition: text-shadow 0.3s ease;
}

.navbar-brand:hover .team-name {
    text-shadow: 0 0 8px rgba(100, 116, 139, 0.4);
}

.acm-icon {
    width: 80px;
    height: 80px;
    margin-left: 8px;
    filter: drop-shadow(0 0 2px rgba(0, 0, 0, 0.2));
}

.auth-buttons {
    display: flex;
    align-items: center;
    gap: 10px;
}

.logout-btn {
    position: relative;
    background: linear-gradient(135deg, #EF4444, #f56c6c);
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 12px;
    font-weight: 600;
    font-size: 1rem;
    border: none;
    cursor: pointer;
    overflow: hidden;
    box-shadow:
        0 4px 6px rgba(239, 68, 68, 0.2),
        0 1px 3px rgba(239, 68, 68, 0.1),
        inset 0 1px 0 rgba(255, 255, 255, 0.15);
    transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.logout-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #f56c6c, #B91C1C);
    opacity: 0;
    transition: opacity 0.4s ease;
    z-index: -1;
}

.logout-btn:hover {
    transform: translateY(-3px) scale(1.02);
    box-shadow:
        0 10px 20px rgba(239, 68, 68, 0.25),
        0 6px 6px rgba(239, 68, 68, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.25);
}

.logout-btn:hover::before {
    opacity: 1;
}

.logout-btn:active {
    transform: translateY(1px) scale(0.98);
    box-shadow:
        0 2px 5px rgba(239, 68, 68, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.hamburger {
    display: none;
    flex-direction: column;
    justify-content: space-around;
    width: 30px;
    height: 25px;
    cursor: pointer;
    z-index: 1001;
    margin-left: 15px;
}

.hamburger .bar {
    height: 3px;
    width: 100%;
    background-color: #1e293b;
    border-radius: 3px;
    transition: all 0.3s ease;
}

/* 移动端样式 */
@media (max-width: 992px) {

    .spacer1,
    .spacer2 {
        flex: 1;
    }
}

@media (max-width: 768px) {
    .navbar {
        padding: 0 1rem;
        height: 60px;
        justify-content: space-between;
    }

    .hamburger {
        display: flex;
        order: 1;
    }

    .navbar-brand {
        order: 1;
    }

    .spacer1,
    .spacer2 {
        display: none;
    }

    /* 隐藏非移动状态下的元素 */
    .nav-links:not(.mobile-active),
    .auth-buttons:not(.mobile-active),
    .login-btn:not(.mobile-active) {
        display: none;
    }

    /* 移动菜单激活状态 */
    .nav-links.mobile-active {
        display: flex;
        position: fixed;
        top: 60px;
        left: 0;
        right: 0;
        background: rgba(248, 250, 252, 0.98);
        backdrop-filter: blur(5px);
        flex-direction: column;
        align-items: center;
        padding: 1rem 0;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
        z-index: 999;
        gap: 0;
    }

    .auth-buttons.mobile-active {
        display: flex;
        position: fixed;
        top: 420px;
        left: 0;
        right: 0;
        background: rgba(248, 250, 252, 0.98);
        flex-direction: column;
        align-items: center;
        padding: 1rem 0;
        gap: 0.5rem;
        z-index: 999;
    }

    .nav-links.mobile-active li {
        width: 100%;
        text-align: center;
        padding: 0.75rem 0;
    }

    .nav-links.mobile-active a,
    .auth-buttons.mobile-active a,
    .auth-buttons.mobile-active button {
        width: 90%;
        display: block;
        margin: 0 auto;
        text-align: center;
        padding: 0.8rem !important;
    }

    .control-panel-btn,
    .logout-btn {
        width: 100%;
        margin: 0.25rem 0;
    }

    .brand-logo {
        height: 40px;
        width: 40px;
    }

    .school-name {
        font-size: 0.75rem;
    }

    .team-name {
        font-size: 0.9375rem;
    }

    .acm-icon {
        width: 50px;
        height: 50px;
    }

    .hamburger.active .bar:nth-child(1) {
        transform: translateY(8px) rotate(45deg);
    }

    .hamburger.active .bar:nth-child(2) {
        opacity: 0;
    }

    .hamburger.active .bar:nth-child(3) {
        transform: translateY(-8px) rotate(-45deg);
    }
}

@media (max-width: 480px) {
    .navbar-brand {
        flex: 1;
    }

    .brand-text {
        flex: 1;
    }

    .login-btn,
    .logout-btn,
    .control-panel-btn {
        padding: 0.6rem 1.2rem !important;
        font-size: 0.9rem;
    }
}
</style>