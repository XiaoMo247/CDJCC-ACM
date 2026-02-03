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
            <el-tag :type="item.type || 'info'" size="small">公告</el-tag>
            <div class="meta-info">
              <span class="view-count">
                <i class="el-icon-view"></i> {{ item.view_count || 0 }} 次查看
              </span>
              <span class="date">{{ item.date }}</span>
            </div>
          </div>
          <h3 class="title">{{ item.title }}</h3>
          <p class="summary">{{ getContentSummary(item.content) }}</p>
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
    <el-dialog
      :lock-scroll="false"
      v-model="dialogVisible"
      :title="currentItem.title"
      width="60%"
      custom-class="announcement-dialog"
    >
      <div class="dialog-content">
        <div class="meta">
          <span class="date">{{ currentItem.date }}</span>
          <span class="view-count">
            <i class="el-icon-view"></i> {{ currentItem.view_count || 0 }} 次查看
          </span>
        </div>
        <div class="markdown-content" v-html="renderedContent" style="max-height: 60vh; overflow-y: auto;"></div>
      </div>
      <template #footer>
        <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

export default {
  setup() {
    const announcements = ref([])
    const loading = ref(true)
    const dialogVisible = ref(false)
    const currentItem = ref({})

    // 配置marked选项
    marked.setOptions({
      breaks: true,      // 支持GFM换行
      gfm: true,         // GitHub Flavored Markdown
      headerIds: false,  // 不生成header ID
      mangle: false      // 不混淆邮箱地址
    })

    const fetchAnnouncements = async () => {
      try {
        const res = await request.get('/announcement/list')
        if (res.data && res.data.announcements) {
          announcements.value = res.data.announcements.map(item => ({
            ...item,
            date: formatDate(item.created_at),
            type: 'info'
          }))
        }
      } catch (err) {
        ElMessage.error('获取公告列表失败')
        console.error(err)
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

    const isHTML = (content) => {
      return content && (content.includes('<p>') || content.includes('<div>') || content.includes('<span>'))
    }

    const renderedContent = computed(() => {
      const content = currentItem.value.content
      if (!content) return ''

      // 智能判断内容格式
      if (isHTML(content)) {
        // 旧的HTML格式公告，直接清理后显示
        return DOMPurify.sanitize(content)
      } else {
        // Markdown格式，先渲染再清理
        const html = marked.parse(content)
        return DOMPurify.sanitize(html)
      }
    })

    const getContentSummary = (content) => {
      if (!content) return '点击查看详情...'

      // 移除Markdown和HTML标记，获取纯文本
      let text = content
        .replace(/<[^>]+>/g, '')           // 移除HTML标签
        .replace(/[#*`>\-\[\]()!]/g, '')  // 移除Markdown标记
        .replace(/\n/g, ' ')               // 替换换行为空格
        .trim()

      return text.length > 80 ? text.substring(0, 80) + '...' : text
    }

    const showDetail = async (item) => {
      currentItem.value = item
      dialogVisible.value = true

      // 增加查看次数
      try {
        const res = await request.post(`/announcement/view/${item.id}`)
        if (res.data && res.data.view_count !== undefined) {
          // 更新本地显示的查看次数
          item.view_count = res.data.view_count
          currentItem.value.view_count = res.data.view_count
        }
      } catch (err) {
        console.error('增加查看次数失败:', err)
        // 不影响用户体验，静默失败
      }
    }

    onMounted(() => {
      fetchAnnouncements()
    })

    return {
      announcements,
      loading,
      dialogVisible,
      currentItem,
      renderedContent,
      getContentSummary,
      showDetail
    }
  }
}
</script>

<style scoped>
.announcements-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  text-align: center;
  margin-bottom: 40px;
}

.header h2 {
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
}

.header h2 i {
  margin-right: 12px;
  color: #409eff;
}

.subtitle {
  font-size: 16px;
  color: #909399;
}

.announcement-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.announcement-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.announcement-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.card-content {
  padding: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.meta-info {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #909399;
}

.view-count {
  display: flex;
  align-items: center;
  gap: 4px;
}

.date {
  color: #909399;
  font-size: 13px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.summary {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.footer {
  display: flex;
  justify-content: flex-end;
}

.detail-btn {
  font-size: 14px;
  padding: 0;
}

/* 对话框样式 */
.dialog-content {
  padding: 20px;
}

.dialog-content .meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  margin-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
  font-size: 14px;
  color: #909399;
}

/* Markdown内容样式 */
.markdown-content {
  line-height: 1.8;
  font-size: 15px;
  color: #303133;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  color: #303133;
}

.markdown-content :deep(h1) {
  font-size: 28px;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 8px;
}

.markdown-content :deep(h2) {
  font-size: 24px;
}

.markdown-content :deep(h3) {
  font-size: 20px;
}

.markdown-content :deep(p) {
  margin-bottom: 16px;
}

.markdown-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 16px 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.markdown-content :deep(code) {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  color: #e03997;
}

.markdown-content :deep(pre) {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 16px 0;
}

.markdown-content :deep(pre code) {
  background: none;
  padding: 0;
  color: #303133;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid #409eff;
  padding-left: 16px;
  margin: 16px 0;
  color: #606266;
  background: #f5f7fa;
  padding: 12px 16px;
  border-radius: 4px;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 24px;
  margin-bottom: 16px;
}

.markdown-content :deep(li) {
  margin-bottom: 8px;
}

.markdown-content :deep(a) {
  color: #409eff;
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.markdown-content :deep(table th),
.markdown-content :deep(table td) {
  border: 1px solid #dcdfe6;
  padding: 8px 12px;
  text-align: left;
}

.markdown-content :deep(table th) {
  background: #f5f7fa;
  font-weight: 600;
}

/* 响应式 */
@media (max-width: 768px) {
  .announcements-container {
    padding: 20px 15px;
  }

  .announcement-list {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .header h2 {
    font-size: 24px;
  }

  :deep(.el-dialog) {
    width: 90% !important;
  }
}
</style>
