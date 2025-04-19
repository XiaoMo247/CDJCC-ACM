<template>
  <div class="join-container">
    <h1 class="title">加入我们</h1>
    <p class="sub-title">如果你对算法竞赛感兴趣，欢迎加入我们！</p>

    <div class="qq-group-box">
      <p class="qq-text">加入官方 QQ 群：</p>
      <p class="qq-number">932140807</p>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.join-container {
  padding: 40px 20px;
  max-width: 600px;
  margin: 0 auto;
  background-color: #f9f9f9;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  animation: fadeInUp 0.5s ease-in-out forwards;
}

.title {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 10px;
  color: #3498db;
}

.sub-title {
  font-size: 1.1rem;
  color: #555;
  margin-bottom: 30px;
}

/* QQ 群信息样式 */
.qq-group-box {
  background-color: #eaf6ff;
  border: 2px dashed #3498db;
  padding: 20px;
  border-radius: 10px;
  margin-bottom: 30px;
}

.qq-text {
  font-size: 1.2rem;
  color: #333;
  margin-bottom: 10px;
}

.qq-number {
  font-size: 1.8rem;
  font-weight: bold;
  color: #3498db;
}

/* 动画效果 */
@keyframes fadeInUp {
  0% {
    opacity: 0;
    transform: translateY(20px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 移动端适配 */
@media (max-width: 500px) {
  .title {
    font-size: 1.5rem;
  }

  .qq-number {
    font-size: 1.4rem;
  }

  .join-container {
    padding: 20px 10px;
  }
}
</style>
