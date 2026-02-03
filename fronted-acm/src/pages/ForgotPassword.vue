<template>
  <div class="forgot-password-page">
    <div class="card">
      <h2>忘记密码？</h2>
      <p>请加入官方 QQ 群并联系管理员重置密码。</p>
      <div class="qq-info">
        <span>QQ群号：</span>
        <strong>932140807</strong>
      </div>
      <a :href="qqGroupUrl" target="_blank" class="join-qq">点我加群</a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ForgotPasswordPage',
  data() {
    return {
      qqGroupUrl: import.meta.env.VITE_QQ_GROUP_URL || 'https://jq.qq.com/?_wv=1027&k=KoVRvLvL'
    }
  }
};
</script>

<style scoped>
.forgot-password-page {
  background: url(../assets/forgot.jpg) no-repeat center center;
  background-size: cover;
  min-height: calc(100vh - 72px - 60px); /* NavBar 72px + Footer 60px */
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  box-sizing: border-box;
}

.card {
  background: linear-gradient(135deg, #f0f4f8 0%, #d9e2ec 100%);
  opacity: 0.95;
  border-radius: 12px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
  padding: 40px;
  text-align: center;
  max-width: 400px;
  width: 100%;
}

.card h2 {
  font-size: 22px;
  margin-bottom: 10px;
  color: #2d3748;
}

.card p {
  font-size: 15px;
  color: #4a5568;
  margin-bottom: 16px;
}

.qq-info {
  font-size: 16px;
  margin-bottom: 20px;
}

.join-qq {
  display: inline-block;
  padding: 10px 20px;
  background-color: #38bdf8;
  color: white;
  border-radius: 6px;
  text-decoration: none;
  transition: all 0.3s;
}

.join-qq:hover {
  background-color: #0ea5e9;
}

/* 移动端适配 */
@media (max-width: 480px) {
  .card {
    padding: 20px;
  }

  .card h2 {
    font-size: 18px;
  }

  .card p,
  .qq-info {
    font-size: 14px;
  }

  .join-qq {
    padding: 8px 16px;
    font-size: 14px;
  }
}
</style>
