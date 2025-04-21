<template>
    <div class="navbar">
        <div class="logo">
            <h1>锦城 ACM</h1>
        </div>
        <div class="spacer"></div>
        <ul class="nav-links">
            <li><router-link to="/">首页</router-link></li>
            <li><router-link to="/knowledge">知识库</router-link></li>
            <li><router-link to="/competition">比赛页面</router-link></li>
            <li><router-link to="/ranking">排行榜</router-link></li>
            <li><router-link to="/faq">FAQ</router-link></li>
            <li><router-link to="/login">登录</router-link></li>
        </ul>
        <div class="spacer"></div>
    </div>
</template>

<script>
export default {
    name: 'NavBar'
}
</script>

<style scoped>
.navbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 60px;
    background-color: #333;
    color: white;
    display: flex;
    align-items: center;
    padding: 0 20px;
    z-index: 1000;
}

.logo h1 {
    margin: 0;
    font-size: 24px;
    white-space: nowrap;
}

.spacer {
    flex: 1;
}

.nav-links {
    list-style: none;
    display: flex;
    gap: 20px;
    margin: 0;
    padding: 0;
}

.nav-links li {
    padding: 5px;
}

.nav-links a {
    color: white;
    text-decoration: none;
    font-weight: bold;
}

.nav-links a:hover {
    text-decoration: underline;
}
</style>