<template>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css" />

  <div class="acm-homepage">
    <!-- 轮播图部分 -->
    <section class="carousel-section">
      <el-carousel :interval="5000" height="500px" arrow="always" indicator-position="outside" trigger="click"
        :autoplay="true">
        <el-carousel-item v-for="(slider, index) in sliders" :key="index">
          <div class="carousel-item">
            <img :src="slider.url" :alt="slider.title" class="carousel-image" />
            <div class="carousel-caption">
              <h3>{{ slider.title }}</h3>
              <p>{{ slider.description }}</p>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
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
import { teamInfo } from './data.js'
import './HomePage.css'

export default {
  setup() {
    const sliders = ref([])
    const loading = ref(false)

    const noticeDialogVisible = ref(false)
    const currentNotice = ref({})

    const showNoticeDetail = (notice) => {
      currentNotice.value = notice
      noticeDialogVisible.value = true
    }

    // 获取轮播图数据
    const fetchSliders = async () => {
      try {
        const res = await request.get('/slider/list')
        if (res.data && res.data.data) {
          sliders.value = res.data.data.map(slider => ({
            url: slider.image,
            title: slider.title,
            description: slider.content,
          }))
        }
      } catch (err) {
        console.error('获取轮播图失败:', err)
        ElMessage.error('获取轮播图列表失败')
      }
    }

    onMounted(() => {
      fetchSliders()
    })

    return {
      teamInfo,
      sliders,
      noticeDialogVisible,
      currentNotice,
      loading,
      showNoticeDetail
    }
  }
}
</script>