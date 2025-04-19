<template>
  <transition name="fade-slide">
    <div class="home-container">
      <!-- 轮播图 -->
      <el-carousel height="360px" autoplay>
        <el-carousel-item v-for="(img, index) in carouselImages" :key="index">
          <div class="carousel-img-container">
            <img :src="img.src" class="carousel-img" />
            <div class="carousel-text">
              <h2>{{ img.caption }}</h2>
              <p>{{ img.subCaption }}</p>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>

      <!-- 公告区域 -->
      <section class="announcement">
        <h2>📢 最新公告</h2>
        <div class="announcement-list">
          <el-card
            v-for="(item, index) in announcements"
            :key="index"
            class="announcement-card"
            @click="showAnnouncement(item)"
            shadow="hover"
          >
            <div class="card-header">
              <h3>{{ item.title }}</h3>
              <span class="date">{{ formatDate(item.created_at) }}</span>
            </div>
            <p class="content-preview">{{ item.content }}</p>
          </el-card>
        </div>
      </section>

      <!-- 简介 -->
      <section class="intro">
        <h2>🎓 关于我们</h2>
        <p>
          成都锦城学院 ACM 集训队致力于培养程序设计能力与算法思维，欢迎对算法、竞赛感兴趣的同学加入我们！
        </p>
        <el-button type="primary" @click="$router.push('/join')">加入我们</el-button>
      </section>

      <!-- 快捷入口 -->
      <section class="quick-links">
        <el-row :gutter="20" justify="center">
          <el-col
            :xs="24" :sm="12" :md="8" :lg="6"
            v-for="item in quickLinks"
            :key="item.path"
          >
            <el-card
              class="quick-card"
              @click="$router.push(item.path)"
              shadow="hover"
            >
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </el-card>
          </el-col>
        </el-row>
      </section>

      <!-- 公告弹窗 -->
      <el-dialog
        v-model="dialogVisible"
        title="公告详情"
        width="80%"
        center
      >
        <h3>{{ selectedAnnouncement.title }}</h3>
        <p class="date">{{ formatDate(selectedAnnouncement.created_at) }}</p>
        <p class="full-content">{{ selectedAnnouncement.content }}</p>
      </el-dialog>
    </div>
  </transition>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 轮播图数据
const carouselImages = [
  { src: '/assets/banner1.jpg', caption: '创新与实践并行', subCaption: '欢迎加入我们，开启算法竞赛之旅！' },
  { src: '/assets/banner2.jpg', caption: '集训队精英汇聚', subCaption: '我们为你提供最优质的训练平台。' },
  { src: '/assets/banner3.jpg', caption: '突破自我挑战极限', subCaption: '成就最强算法队伍，等待你的加入。' },
]

const announcements = ref([])
const selectedAnnouncement = ref({})
const dialogVisible = ref(false)

const quickLinks = [
  { title: '比赛信息', desc: '查看即将到来的比赛', path: '/contests' },
  { title: '排行榜', desc: '查看队员训练排名', path: '/rank' },
  { title: '加入我们', desc: '想加入ACM集训队？点我报名', path: '/join' },
]

const getAnnouncements = async () => {
  try {
    const res = await axios.get('http://localhost:8081/api/announcement/list', {
      headers: { Authorization: `Bearer ${localStorage.getItem('adminToken')}` },
    })
    announcements.value = res.data.announcements || []
  } catch (error) {
    ElMessage.error('获取公告失败')
  }
}

const formatDate = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

const showAnnouncement = (item) => {
  selectedAnnouncement.value = item
  dialogVisible.value = true
}

onMounted(() => getAnnouncements())

</script>

<style scoped>
.home-container {
  padding: 40px 20px;
}

/* 动效过渡 */
.fade-slide-enter-active, .fade-slide-leave-active {
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

/* 轮播图 */
.carousel-img-container {
  position: relative;
}
.carousel-img {
  width: 100%;
  height: 360px;
  object-fit: cover;
}
.carousel-text {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.5);
  padding: 16px 24px;
  color: #fff;
  text-align: center;
  border-radius: 8px;
}

/* 公告 */
.announcement {
  margin-top: 40px;
}
.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.announcement-card {
  transition: transform 0.3s ease;
  cursor: pointer;
}
.announcement-card:hover {
  transform: translateY(-3px);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.content-preview {
  color: #4b5563;
  margin-top: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.full-content {
  white-space: pre-line;
}

/* 简介 */
.intro {
  margin: 60px auto 40px;
  text-align: center;
}
.intro p {
  margin: 16px 0;
}

/* 快捷入口 */
.quick-links {
  margin: 40px 0;
}
.quick-card {
  text-align: center;
  cursor: pointer;
  transition: transform 0.3s ease;
}
.quick-card:hover {
  transform: translateY(-4px);
}

/* 响应式支持 */
@media (max-width: 768px) {
  .carousel-text {
    font-size: 16px;
    padding: 12px;
  }

  .intro p,
  .full-content {
    font-size: 14px;
  }

  .announcement-list {
    gap: 12px;
  }
}
</style>
