<template>
  <footer class="acm-footer">
    <div class="footer-content">
      <p>© {{ year }} 成都锦城学院 ACM 集训队 · 保留所有权利</p>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>

<style scoped>
.acm-footer {
  background-color: #0f172a;
  color: #94a3b8;
  text-align: center;
  padding: 20px 16px;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.5px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  margin-top: auto;
}

.footer-content {
  max-width: 1280px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
}

.footer-content p {
  margin: 0;
  color: #94a3b8;
  transition: color 0.3s;
}

.footer-content p:hover {
  color: #38bdf8;
}
</style>
