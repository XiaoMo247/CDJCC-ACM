<template>
  <div class="announcements-container page-announcement">
    <!-- 标题区 -->
    <div class="header page-header page-hero">
      <h2><i class="fas fa-bell"></i> 团队公告</h2>
      <p class="subtitle">最新动态与重要通知</p>
    </div>

    <!-- 加载状态 -->
    <el-skeleton v-if="loading" :rows="5" animated />

    <!-- 公告列表 -->
    <div class="announcement-list" v-else>
      <el-card 
        v-for="(item, index) in announcements" 
        :key="index"
        class="announcement-card"
        shadow="hover"
        @click="showDetail(item)"
      >
        <div class="card-content">
          <div class="card-header">
            <el-tag :type="item.type || 'info'" size="small">{{ item.tag || '公告' }}</el-tag>
            <span class="date">{{ item.date }}</span>
          </div>
          <h3 class="title">{{ item.title }}</h3>
          <p class="summary">{{ item.summary || '点击查看详情...' }}</p>
          <div class="footer">
            <el-button type="text" class="detail-btn">
              查看详情 <i class="el-icon-arrow-right"></i>
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <el-empty v-if="!loading && announcements.length === 0" description="暂无公告" />

    <!-- 详情对话框 -->
    <el-dialog :lock-scroll="false" 
      v-model="dialogVisible" 
      :title="currentItem.title" 
      width="50%"
      custom-class="announcement-dialog"
    >
      <div class="dialog-content">
        <div class="meta">
          <span class="date">{{ currentItem.date }}</span>
          <el-tag v-if="currentItem.tag" :type="currentItem.type">{{ currentItem.tag }}</el-tag>
        </div>
        <div class="content" v-html="sanitizedContent(currentItem.content)" style="max-height: 60vh; overflow-y: auto;"></div>
      </div>
      <template #footer>
        <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';
import request from '@/utils/request';
import { ElMessage } from 'element-plus';
import DOMPurify from 'dompurify';

export default {
  setup() {
    const announcements = ref([]);
    const loading = ref(true);
    const dialogVisible = ref(false);
    const currentItem = ref({});

    const fetchAnnouncements = async () => {
      try {
        const res = await request.get('/announcement/list');
        if (res.data && res.data.announcements) {
          announcements.value = res.data.announcements.map(item => ({
            ...item,
            date: formatDate(item.created_at),
            type: item.tagType || 'info',
            content: item.detailContent || item.content
          }));
        }
      } catch (err) {
        ElMessage.error('获取公告列表失败');
        console.error(err);
      } finally {
        loading.value = false;
      }
    };

    const formatDate = (timestamp) => {
      const date = new Date(timestamp);
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      });
    };

    const sanitizedContent = (html) => {
      return DOMPurify.sanitize(html);
    };

    const showDetail = (item) => {
      currentItem.value = item;
      dialogVisible.value = true;
    };

    const abortController = new AbortController();

    onMounted(() => {
      fetchAnnouncements();
    });

    onUnmounted(() => {
      abortController.abort();
    });

    return {
      announcements,
      loading,
      dialogVisible,
      currentItem,
      showDetail,
      sanitizedContent
    };
  }
};
</script>

<style>
body {
  overflow-y: scroll !important;
}
</style>

<style scoped>
/* 原有样式保持不变 */
.announcements-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  background: linear-gradient(135deg, #a5f3fc 0%, #7dd3fc 100%);
  color: #0c4a6e;
  text-align: center;
  margin-bottom: 60px;
  padding: 60px 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.header h2 {
  font-size: 2.5rem;
  font-weight: 700;
  color: white;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
}

.header h2 i {
  font-size: 2.8rem;
}

.subtitle {
  color: rgba(255, 255, 255, 0.95);
  margin-top: 10px;
  font-size: 1.125rem;
}

.announcement-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 32px;
  padding: 0;
}

.announcement-card {
  cursor: pointer;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 12px;
  position: relative;
  overflow: hidden;
  border-left: 3px solid #a5f3fc;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.announcement-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(135deg, rgba(165, 243, 252, 0.02) 0%, transparent 100%);
  pointer-events: none;
}

.announcement-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(165, 243, 252, 0.1);
  border-left-color: #7dd3fc;
}

.card-content {
  padding: 24px;
  position: relative;
  z-index: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.date {
  color: #999;
  font-size: 14px;
}

.title {
  font-size: 18px;
  margin: 10px 0;
  color: #333;
}

.summary {
  color: #666;
  line-height: 1.6;
  margin-bottom: 15px;
}

.footer {
  text-align: right;
}

/* 对话框样式 */
.announcement-dialog {
  border-radius: 12px;
  max-width: 80vw;
  width: 800px !important;
}

.announcement-dialog .el-dialog__body {
  padding: 20px;
}

.dialog-content {
  padding: 0 20px;
}

.meta {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.meta .date {
  margin-right: 15px;
}

.content {
  line-height: 1.8;
}
</style>
