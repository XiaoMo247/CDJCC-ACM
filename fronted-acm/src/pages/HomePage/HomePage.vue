<template>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css" />

  <div class="acm-homepage">
    <!-- 轮播图部分 -->
    <section class="carousel-section">
      <el-carousel :interval="5000" height="500px" arrow="always" indicator-position="outside" trigger="click"
        :autoplay="true">
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
    </section>

    <!-- 公告栏部分 -->
    <section class="notice-section">
      <div class="section-header">
        <h2 class="section-title">
          <i class="fas fa-bell title-icon"></i>
          最新公告
        </h2>
        <div class="section-subtitle">及时了解团队动态和训练安排</div>
      </div>

      <el-card class="notice-card" shadow="hover">
        <el-scrollbar height="400px">
          <el-timeline>
            <el-timeline-item v-for="(notice, index) in notices" :key="index" :timestamp="notice.date" placement="top"
              :type="notice.type" :color="notice.color" :hollow="index === 0" :size="index === 0 ? 'large' : 'normal'">
              <el-card class="timeline-card" shadow="never">
                <h3>{{ notice.title }}</h3>
                <p class="notice-content">{{ notice.content }}</p>
                <div class="notice-footer">
                  <el-tag v-if="notice.tag" :type="notice.tagType">{{ notice.tag }}</el-tag>
                  <el-button v-if="notice.hasDetail" type="text" class="detail-btn" @click="showNoticeDetail(notice)">
                    查看详情 <i class="fas fa-arrow-right"></i>
                  </el-button>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-scrollbar>
      </el-card>
    </section>

    <!-- 团队介绍部分 -->
    <section class="team-section">
      <div class="section-header">
        <h2 class="section-title">
          <i class="fas fa-users title-icon"></i>
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
                <span>所属学院: {{ teamInfo.college }}  |</span>
                <span>成立时间: {{ teamInfo.established }}  |</span>
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
                <i :class="['feature-icon', feature.icon]"></i>
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

    <!-- 荣誉墙部分 -->
    <section class="honor-section">
      <div class="section-header">
        <h2 class="section-title">
          <i class="fas fa-trophy title-icon"></i>
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
          <div class="honor-item" v-for="(honor, index) in filteredHonors" :key="index"
            :class="'honor-level-' + honor.level">
            <div class="honor-badge">
              <i v-if="honor.level === 'gold'" class="fas fa-medal"></i>
              <i v-else-if="honor.level === 'silver'" class="fas fa-coins"></i>
              <i v-else class="fas fa-crosshairs"></i>
            </div>
            <div class="honor-content">
              <h3>{{ honor.title }}</h3>
              <div class="honor-meta">
                <span class="honor-year">{{ honor.year }}</span>
                <el-tag :type="getHonorTagType(honor.type)" size="small">{{ honor.typeText }}</el-tag>
              </div>
              <p class="honor-desc">{{ honor.description }}</p>
              <div class="honor-members">
                <el-tooltip v-for="member in honor.members" :key="member.name"
                  :content="member.name + (member.class ? ' (' + member.class + ')' : '')" placement="top">
                  <el-avatar :size="32" :src="member.avatar" />
                </el-tooltip>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </section>

    <!-- 公告详情对话框 -->
    <el-dialog v-model="noticeDialogVisible" :title="currentNotice.title" width="60%" center
      :close-on-click-modal="false">
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

<script>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { teamInfo, photos, honors } from './data.js'
import './HomePage.css'

export default {
  setup() {
    const notices = ref([])
    const loading = ref(false)

    const activeHonorTab = ref('all')
    const filteredHonors = computed(() => {
      if (activeHonorTab.value === 'all') return honors
      return honors.filter(honor => honor.type === activeHonorTab.value)
    })

    const noticeDialogVisible = ref(false)
    const currentNotice = ref({})

    const getHonorTagType = (type) => {
      switch (type) {
        case 'icpc':
          return 'success'
        case 'ccpc':
          return 'warning'
        default:
          return 'info'
      }
    }

    const showNoticeDetail = (notice) => {
      currentNotice.value = notice
      noticeDialogVisible.value = true
    }

    const fetchAnnouncements = async () => {
      loading.value = true
      try {
        const res = await request.get('/announcement/list')
        if (res.data && res.data.announcements) {
          notices.value = res.data.announcements.map(item => ({
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
        loading.value = false
      }
    }

    const formatDate = (timestamp) => {
      const date = new Date(timestamp)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    onMounted(() => {
      fetchAnnouncements()
    })

    return {
      teamInfo,
      photos,
      notices,
      filteredHonors,
      activeHonorTab,
      noticeDialogVisible,
      currentNotice,
      loading,
      getHonorTagType,
      showNoticeDetail
    }
  }
}
</script>
