<template>
    <div class="acm-homepage">
      <!-- 主内容区 -->
      <main class="acm-main">
        <!-- 公告栏 -->
        <section class="notice-section">
          <div class="section-header">
            <h2 class="section-title">
              <el-icon class="title-icon"><Bell /></el-icon>
              最新公告
            </h2>
            <div class="section-subtitle">及时了解团队动态和训练安排</div>
          </div>
          
          <el-card class="notice-card" shadow="hover">
            <el-scrollbar height="400px">
              <el-timeline>
                <el-timeline-item
                  v-for="(notice, index) in notices"
                  :key="index"
                  :timestamp="notice.date"
                  placement="top"
                  :type="notice.type"
                  :color="notice.color"
                  :hollow="index === 0"
                  :size="index === 0 ? 'large' : 'normal'"
                >
                  <el-card class="timeline-card" shadow="never">
                    <h3>{{ notice.title }}</h3>
                    <p class="notice-content">{{ notice.content }}</p>
                    <div class="notice-footer">
                      <el-tag v-if="notice.tag" :type="notice.tagType">{{ notice.tag }}</el-tag>
                      <el-button
                        v-if="notice.hasDetail"
                        type="text"
                        class="detail-btn"
                        @click="showNoticeDetail(notice)"
                      >
                        查看详情 <el-icon><ArrowRight /></el-icon>
                      </el-button>
                    </div>
                  </el-card>
                </el-timeline-item>
              </el-timeline>
            </el-scrollbar>
          </el-card>
        </section>
  
        <!-- 团队介绍 -->
        <section class="team-section">
          <div class="section-header">
            <h2 class="section-title">
              <el-icon class="title-icon"><User /></el-icon>
              团队简介
            </h2>
            <div class="section-subtitle">我们的使命与愿景</div>
          </div>
          
          <el-card class="team-card" shadow="hover">
            <div class="team-content">
              <div class="team-avatar">
                <el-avatar :size="120" :src="teamInfo.avatar" />
                <div class="team-name">{{ teamInfo.name }}</div>
              </div>
              <div class="team-details">
                <div class="team-meta">
                  <div class="meta-item">
                    <el-icon><School /></el-icon>
                    <span>隶属学院: {{ teamInfo.college }}</span>
                  </div>
                  <div class="meta-item">
                    <el-icon><Flag /></el-icon>
                    <span>成立时间: {{ teamInfo.established }}</span>
                  </div>
                  <div class="meta-item">
                    <el-icon><Trophy /></el-icon>
                    <span>获奖总数: {{ teamInfo.awards }}项</span>
                  </div>
                </div>
                
                <el-divider />
                
                <div class="team-description">
                  <p>{{ teamInfo.description }}</p>
                  <p>{{ teamInfo.mission }}</p>
                </div>
                
                <el-divider />
                
                <div class="team-features">
                  <div class="feature-item" v-for="(feature, index) in teamInfo.features" :key="index">
                    <el-icon class="feature-icon"><component :is="feature.icon" /></el-icon>
                    <div class="feature-content">
                      <h4>{{ feature.title }}</h4>
                      <p>{{ feature.desc }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </section>
  
        <!-- 荣誉墙 -->
        <section class="honor-section">
          <div class="section-header">
            <h2 class="section-title">
              <el-icon class="title-icon"><Trophy /></el-icon>
              荣誉墙
            </h2>
            <div class="section-subtitle">我们的辉煌成就</div>
          </div>
          
          <el-card class="honor-card" shadow="hover">
            <div class="honor-tabs">
              <el-tabs v-model="activeHonorTab" type="card">
                <el-tab-pane label="全部荣誉" name="all"></el-tab-pane>
                <el-tab-pane label="ICPC" name="icpc"></el-tab-pane>
                <el-tab-pane label="CCPC" name="ccpc"></el-tab-pane>
                <el-tab-pane label="省级比赛" name="provincial"></el-tab-pane>
              </el-tabs>
            </div>
            
            <div class="honor-grid">
              <div
                class="honor-item"
                v-for="(honor, index) in filteredHonors"
                :key="index"
                :class="'honor-level-' + honor.level"
              >
                <div class="honor-badge">
                  <el-icon v-if="honor.level === 'gold'"><Medal /></el-icon>
                  <el-icon v-else-if="honor.level === 'silver'"><Coin /></el-icon>
                  <el-icon v-else><Aim /></el-icon>
                </div>
                <div class="honor-content">
                  <h3>{{ honor.title }}</h3>
                  <div class="honor-meta">
                    <span class="honor-year">{{ honor.year }}</span>
                    <el-tag :type="getHonorTagType(honor.type)" size="small">{{ honor.typeText }}</el-tag>
                  </div>
                  <p class="honor-desc">{{ honor.description }}</p>
                  <div class="honor-members">
                    <el-tooltip
                      v-for="member in honor.members"
                      :key="member.name"
                      :content="member.name + (member.class ? ' (' + member.class + ')' : '')"
                      placement="top"
                    >
                      <el-avatar :size="32" :src="member.avatar" />
                    </el-tooltip>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </section>
  
        <!-- 比赛风采 -->
        <section class="gallery-section">
          <div class="section-header">
            <h2 class="section-title">
              <el-icon class="title-icon"><Camera /></el-icon>
              比赛风采
            </h2>
            <div class="section-subtitle">精彩瞬间，永恒记忆</div>
          </div>
          
          <el-card class="gallery-card" shadow="hover">
            <el-carousel
              :interval="5000"
              height="500px"
              arrow="always"
              indicator-position="outside"
              trigger="click"
              :autoplay="true"
            >
              <el-carousel-item v-for="(photo, index) in photos" :key="index">
                <div class="carousel-item">
                  <img :src="photo.url" :alt="photo.title" class="carousel-image" />
                  <div class="carousel-caption">
                    <h3>{{ photo.title }}</h3>
                    <p>{{ photo.description }}</p>
                    <div class="photo-meta">
                      <span class="photo-date">{{ photo.date }}</span>
                      <span class="photo-location">{{ photo.location }}</span>
                    </div>
                  </div>
                </div>
              </el-carousel-item>
            </el-carousel>
            
            <div class="thumbnail-container">
              <div
                class="thumbnail-item"
                v-for="(photo, index) in photos"
                :key="'thumb-' + index"
                @click="switchPhoto(index)"
                :class="{ active: activePhotoIndex === index }"
              >
                <img :src="photo.url" :alt="photo.title" />
              </div>
            </div>
          </el-card>
        </section>
      </main>
  
      <!-- 公告详情对话框 -->
      <el-dialog
        v-model="noticeDialogVisible"
        :title="currentNotice.title"
        width="60%"
        center
        :close-on-click-modal="false"
      >
        <div class="notice-dialog-content">
          <div class="notice-meta">
            <span class="notice-date">{{ currentNotice.date }}</span>
            <el-tag v-if="currentNotice.tag" :type="currentNotice.tagType">{{ currentNotice.tag }}</el-tag>
          </div>
          <div class="notice-body" v-html="currentNotice.detailContent"></div>
        </div>
        <template #footer>
          <el-button type="primary" @click="noticeDialogVisible = false">关闭</el-button>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Bell, User, Trophy, Camera, School, Flag, Medal, Coin, Aim } from '@element-plus/icons-vue'
import request from '@/utils/request'

// 类型定义
interface Announcement {
  id: number
  title: string
  content: string
  created_at: string
  updated_at: string
  date?: string
  type?: string
  color?: string
  tag?: string
  tagType?: string
  hasDetail?: boolean
  detailContent?: string
}

interface Honor {
  level: string
  title: string
  year: string
  type: string
  typeText: string
  description: string
  members: Array<{
    name: string
    avatar: string
    class?: string
  }>
}

interface Photo {
  url: string
  title: string
  description: string
  date: string
  location: string
}

// 状态管理
const notices = ref<Announcement[]>([])
const loading = ref({
  announcements: false,
  honors: false,
  photos: false
})

// 团队信息配置
const teamInfo = ref({
  name: 'ACM竞赛团队',
  college: '成都锦城学院计算机与软件学院',
  established: '2020年',
  avatar: '/src/assets/acm-logo.jpg',
  description: '本团队致力于算法学习与ACM竞赛训练，培养团队协作精神与编程能力。',
  awards: 28,
  mission: '培养优秀的程序设计人才，在竞赛中展现锦城风采。',
  features: [
    {
      icon: 'User',
      title: '专业指导',
      desc: '由资深算法竞赛指导老师带队'
    },
    {
      icon: 'Clock',
      title: '系统训练',
      desc: '每周固定训练时间，循序渐进提升能力'
    },
    {
      icon: 'Trophy',
      title: '竞赛实战',
      desc: '定期参加各类高水平程序设计竞赛'
    },
    {
      icon: 'Promotion',
      title: '成长空间',
      desc: '提供广阔的技术发展和就业机会'
    }
  ]
})

// 荣誉墙相关
const activeHonorTab = ref('all')
const filteredHonors = ref<Honor[]>([])

// 照片展示相关
const photos = ref<Photo[]>([])
const activePhotoIndex = ref(0)

// 公告详情相关
const noticeDialogVisible = ref(false)
const currentNotice = ref<Announcement>({} as Announcement)

// 方法定义
const getHonorTagType = (type: string) => {
  switch (type) {
    case 'icpc':
      return 'success'
    case 'ccpc':
      return 'warning'
    default:
      return 'info'
  }
}

const switchPhoto = (index: number) => {
  activePhotoIndex.value = index
}

// 获取公告列表
const fetchAnnouncements = async () => {
  loading.value.announcements = true
  try {
    const res = await request.get('/announcement/list')
    if (res.data && res.data.announcements) {
      notices.value = res.data.announcements.map((item: Announcement) => ({
        ...item,
        date: formatDate(item.created_at),
        type: 'primary',
        tag: '公告',
        tagType: 'info'
      }))
    }
  } catch (err) {
    console.error('获取公告失败:', err)
    ElMessage.error('获取公告列表失败，请稍后重试')
  } finally {
    loading.value.announcements = false
  }
}

// 格式化日期
const formatDate = (timestamp: string) => {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const showNoticeDetail = (notice: Announcement & { date?: string }) => {
  currentNotice.value = notice
  noticeDialogVisible.value = true
}

// 生命周期钩子
onMounted(() => {
  fetchAnnouncements()
})
</script>
  
  <style scoped>
  /* 基础样式 */
  .acm-homepage {
    font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
    color: #333;
    line-height: 1.6;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }
  
  /* 通用部分样式 */
  .section-header {
    margin-bottom: 24px;
    text-align: center;
  }
  
  .section-title {
    font-size: 28px;
    font-weight: 600;
    color: #2c3e50;
    margin-bottom: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .title-icon {
    margin-right: 10px;
  }
  
  .section-subtitle {
    font-size: 16px;
    color: #7f8c8d;
  }
  
  /* 公告栏样式 */
  .notice-section {
    margin-bottom: 40px;
  }
  
  .notice-card {
    border-radius: 12px;
    overflow: hidden;
  }
  
  .notice-content {
    margin: 10px 0;
    color: #555;
  }
  
  .notice-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 15px;
  }
  
  .detail-btn {
    padding-right: 0;
  }
  
  .timeline-card {
    border: none;
    background-color: #f8f9fa;
    border-radius: 8px;
    transition: all 0.3s ease;
  }
  
  .timeline-card:hover {
    background-color: #f1f5f9;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }
  
  /* 团队介绍样式 */
  .team-section {
    margin-bottom: 40px;
  }
  
  .team-card {
    border-radius: 12px;
  }
  
  .team-content {
    display: flex;
    flex-direction: row;
    padding: 20px;
  }
  
  .team-avatar {
    flex: 0 0 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-right: 30px;
    border-right: 1px solid #eee;
  }
  
  .team-name {
    margin-top: 15px;
    font-size: 20px;
    font-weight: 600;
    color: #2c3e50;
  }
  
  .team-details {
    flex: 1;
    padding-left: 30px;
  }
  
  .team-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    margin-bottom: 15px;
  }
  
  .meta-item {
    display: flex;
    align-items: center;
    font-size: 15px;
    color: #555;
  }
  
  .meta-item .el-icon {
    margin-right: 8px;
    color: #409EFF;
  }
  
  .team-description {
    margin: 20px 0;
    color: #555;
    line-height: 1.8;
  }
  
  .team-features {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    margin-top: 20px;
  }
  
  .feature-item {
    display: flex;
    padding: 15px;
    background-color: #f8f9fa;
    border-radius: 8px;
    transition: all 0.3s ease;
  }
  
  .feature-item:hover {
    background-color: #f1f5f9;
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }
  
  .feature-icon {
    font-size: 24px;
    color: #409EFF;
    margin-right: 15px;
    flex: 0 0 40px;
  }
  
  .feature-content h4 {
    margin: 0 0 8px 0;
    font-size: 16px;
    color: #2c3e50;
  }
  
  .feature-content p {
    margin: 0;
    font-size: 14px;
    color: #666;
  }
  
  /* 荣誉墙样式 */
  .honor-section {
    margin-bottom: 40px;
  }
  
  .honor-card {
    border-radius: 12px;
    overflow: hidden;
  }
  
  .honor-tabs {
    margin-bottom: 20px;
  }
  
  .honor-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 20px;
  }
  
  .honor-item {
    display: flex;
    padding: 20px;
    border-radius: 8px;
    background-color: #f8f9fa;
    transition: all 0.3s ease;
  }
  
  .honor-item:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .honor-level-gold {
    border-left: 4px solid #E6A23C;
  }
  
  .honor-level-silver {
    border-left: 4px solid #909399;
  }
  
  .honor-level-bronze {
    border-left: 4px solid #F56C6C;
  }
  
  .honor-badge {
    flex: 0 0 50px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    margin-right: 15px;
  }
  
  .honor-level-gold .honor-badge {
    color: #E6A23C;
  }
  
  .honor-level-silver .honor-badge {
    color: #909399;
  }
  
  .honor-level-bronze .honor-badge {
    color: #F56C6C;
  }
  
  .honor-content {
    flex: 1;
  }
  
  .honor-content h3 {
    margin: 0 0 10px 0;
    font-size: 18px;
    color: #2c3e50;
  }
  
  .honor-meta {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    gap: 10px;
  }
  
  .honor-year {
    font-size: 14px;
    color: #7f8c8d;
  }
  
  .honor-desc {
    margin: 10px 0;
    font-size: 14px;
    color: #555;
    line-height: 1.6;
  }
  
  .honor-members {
    display: flex;
    gap: 8px;
    margin-top: 15px;
  }
  
  /* 比赛风采样式 */
  .gallery-section {
    margin-bottom: 40px;
  }
  
  .gallery-card {
    border-radius: 12px;
    overflow: hidden;
  }
  
  .carousel-item {
    position: relative;
    height: 100%;
  }
  
  .carousel-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .carousel-caption {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 30px;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
    color: white;
  }
  
  .carousel-caption h3 {
    margin: 0 0 10px 0;
    font-size: 24px;
  }
  
  .carousel-caption p {
    margin: 0 0 10px 0;
    font-size: 16px;
  }
  
  .photo-meta {
    display: flex;
    gap: 15px;
    font-size: 14px;
    opacity: 0.8;
  }
  
  .thumbnail-container {
    display: flex;
    gap: 10px;
    padding: 15px;
    overflow-x: auto;
  }
  
  .thumbnail-item {
    flex: 0 0 80px;
    height: 60px;
    border-radius: 4px;
    overflow: hidden;
    cursor: pointer;
    opacity: 0.6;
    transition: all 0.3s ease;
  }
  
  .thumbnail-item:hover, .thumbnail-item.active {
    opacity: 1;
    transform: scale(1.05);
  }
  
  .thumbnail-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  /* 公告对话框样式 */
  .notice-dialog-content {
    padding: 0 20px;
  }
  
  .notice-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
    padding-bottom: 10px;
    border-bottom: 1px solid #eee;
  }
  
  .notice-date {
    color: #7f8c8d;
    font-size: 14px;
  }
  
  .notice-body {
    line-height: 1.8;
  }
  
  .notice-body h3 {
    margin-top: 20px;
    color: #2c3e50;
  }
  
  .notice-body ul {
    padding-left: 20px;
  }
  
  .notice-body li {
    margin-bottom: 8px;
  }
  
  /* 响应式设计 */
  @media (max-width: 992px) {
    .team-content {
      flex-direction: column;
    }
    
    .team-avatar {
      flex: auto;
      padding-right: 0;
      padding-bottom: 20px;
      border-right: none;
      border-bottom: 1px solid #eee;
      margin-bottom: 20px;
    }
    
    .team-details {
      padding-left: 0;
    }
    
    .team-features {
      grid-template-columns: 1fr;
    }
    
    .honor-grid {
      grid-template-columns: 1fr;
    }
  }
  
  @media (max-width: 768px) {
    .section-title {
      font-size: 24px;
    }
    
    .carousel-caption {
      padding: 15px;
    }
    
    .carousel-caption h3 {
      font-size: 20px;
    }
    
    .carousel-caption p {
      font-size: 14px;
    }
    
    .thumbnail-item {
      flex: 0 0 60px;
      height: 45px;
    }
  }
  
  @media (max-width: 576px) {
    .acm-homepage {
      padding: 15px;
    }
    
    .section-title {
      font-size: 22px;
    }
    
    .team-meta {
      flex-direction: column;
      gap: 10px;
    }
    
    .carousel-caption {
      padding: 10px;
    }
    
    .carousel-caption h3 {
      font-size: 18px;
      margin-bottom: 5px;
    }
    
    .carousel-caption p {
      display: none;
    }
    
    .photo-meta {
      font-size: 12px;
      gap: 8px;
    }
  }
  
  /* 动画效果 */
  .timeline-card {
    transition: all 0.3s ease;
  }
  
  .timeline-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .honor-item {
    transition: all 0.3s ease;
  }
  
  .honor-item:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .feature-item {
    transition: all 0.3s ease;
  }
  
  .feature-item:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .thumbnail-item {
    transition: all 0.3s ease;
  }
  
  /* 滚动条样式 */
  .el-scrollbar__bar.is-vertical {
    right: 2px;
  }
  
  .el-scrollbar__thumb {
    background-color: rgba(144, 147, 153, 0.3);
  }
  
  .el-scrollbar__thumb:hover {
    background-color: rgba(144, 147, 153, 0.5);
  }
  </style>