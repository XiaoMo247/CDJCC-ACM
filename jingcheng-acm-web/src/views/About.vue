<template>
  <transition name="fade-slide">
    <div class="about-container">
      <h1 class="page-title">关于我们</h1>

      <section class="team-intro">
        <h2 class="section-title">💡 团队简介</h2>
        <p class="intro-text">
          成都锦城学院 ACM 集训队成立于 20XX 年，致力于为对算法和程序设计感兴趣的同学提供一个优质的训练平台。
          我们以“提高编程能力、突破思维局限”为理念，定期组织各种形式的训练与比赛，助力队员在 ICPC、CCPC 等赛事中脱颖而出。
        </p>
      </section>
    </div>
  </transition>
</template>

<script setup>
</script>

<style scoped>
.about-container {
  padding: 40px 20px;
  max-width: 1000px;
  margin: auto;
}

.page-title {
  text-align: center;
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 40px;
  color: #333;
}

.section-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #2c3e50;
  border-left: 4px solid #409EFF;
  padding-left: 12px;
}

.intro-text {
  font-size: 16px;
  line-height: 1.8;
  color: #555;
}

/* 页面切换动效 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.5s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-30px);
}

/* 响应式适配 */
@media (max-width: 600px) {
  .page-title {
    font-size: 28px;
  }
}
</style>