<template>
  <div class="faq-page">
    <div class="header">
      <h1 class="title">✨ 常见问题（FAQ）</h1>
      <p class="subtitle">这里整理了大家最关心的问题和解答！</p>
    </div>

    <div class="search-box">
      <el-input 
        v-model="searchText"
        placeholder="输入关键词快速查找问题..."
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <transition name="fade-stagger" mode="out-in">
      <div v-if="filteredFaqs.length > 0" key="content">
        <el-collapse 
          v-model="activeNames" 
          accordion 
          class="faq-list"
        >
          <transition-group name="list">
            <el-collapse-item
              v-for="(faq, index) in filteredFaqs"
              :key="faq.id"
              :name="index.toString()"
              class="faq-item"
            >
              <template #title>
                <div class="faq-header">
                  <span class="question-icon">❓</span>
                  <span class="faq-question">{{ faq.question }}</span>
                </div>
              </template>
              <div class="faq-answer">
                <span class="answer-icon">💡</span>
                <div class="answer-content">{{ faq.answer }}</div>
              </div>
            </el-collapse-item>
          </transition-group>
        </el-collapse>
      </div>
      
      <div v-else key="empty" class="empty-state">
        <el-empty description="暂时没有找到相关问题">
          <el-button type="primary" @click="resetSearch">重置搜索</el-button>
        </el-empty>
      </div>
    </transition>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import request from '@/utils/request'
import { Search } from '@element-plus/icons-vue'

export default {
  name: 'FaqPage',
  components: { Search },
  setup() {
    const faqList = ref([])
    const activeNames = ref('0')
    const searchText = ref('')

    const filteredFaqs = computed(() => {
      return faqList.value.filter(item => 
        item.question.includes(searchText.value) ||
        item.answer.includes(searchText.value))
    })

    const getFaqList = async () => {
      try {
        const res = await request.get('/faq/list')
        faqList.value = res.data.data || []
      } catch (error) {
        console.error('获取FAQ失败', error)
      }
    }

    const resetSearch = () => {
      searchText.value = ''
    }

    onMounted(() => {
      getFaqList()
    })

    return {
      faqList,
      activeNames,
      searchText,
      filteredFaqs,
      resetSearch
    }
  }
}
</script>

<style scoped>
.faq-page {
  padding: 2rem 1.5rem;
  max-width: 900px;
  margin: 0 auto;
  min-height: 80vh;
}

.header {
  text-align: center;
  margin-bottom: 2.5rem;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  padding: 3rem 1rem;
  border-radius: 1.5rem;
  box-shadow: 0 10px 30px rgba(79, 172, 254, 0.15);
  color: white;
}

.title {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 1rem;
  text-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.subtitle {
  font-size: 1.25rem;
  opacity: 0.95;
  font-weight: 300;
}

.search-box {
  max-width: 600px;
  margin: 0 auto 3rem;
}

.search-input {
  border-radius: 50px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 50px;
  padding: 1rem 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.faq-list {
  background: white;
  border-radius: 1.5rem;
  box-shadow: 0 8px 30px rgba(0,0,0,0.05);
  overflow: hidden;
}

.faq-item {
  transition: all 0.3s ease;
  border-bottom: 1px solid #f0f4f8;
}

.faq-item:last-child {
  border-bottom: none;
}

.faq-header {
  display: flex;
  align-items: center;
  padding: 1rem 0;
}

.question-icon {
  font-size: 1.5rem;
  margin-right: 1rem;
  flex-shrink: 0;
}

.faq-question {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2c3e50;
}

.faq-answer {
  display: flex;
  align-items: flex-start;
  padding: 1.5rem;
  background: #f9fbfd;
  margin: 0 -1rem;
  border-radius: 0 0 1rem 1rem;
}

.answer-icon {
  font-size: 1.2rem;
  margin-right: 1rem;
  flex-shrink: 0;
}

.answer-content {
  font-size: 1rem;
  color: #4a5568;
  line-height: 1.8;
  flex-grow: 1;
}

.empty-state {
  text-align: center;
  padding: 4rem 0;
  background: white;
  border-radius: 1.5rem;
  box-shadow: 0 8px 30px rgba(0,0,0,0.05);
}

/* 动画效果 */
.fade-stagger-enter-active {
  transition: all 0.3s ease;
}

.fade-stagger-leave-active {
  transition: all 0.2s ease;
}

.fade-stagger-enter-from,
.fade-stagger-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.list-move {
  transition: transform 0.3s ease;
}

@media (max-width: 768px) {
  .title {
    font-size: 2rem;
  }
  
  .faq-question {
    font-size: 1rem;
  }
  
  .answer-content {
    font-size: 0.9rem;
  }
  
  .faq-answer {
    padding: 1rem;
  }
}
</style>