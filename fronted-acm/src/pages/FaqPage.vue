<template>
  <div class="faq-page">
    <div class="header">
      <h1 class="title">💡常见问题</h1>
      <p class="subtitle">传递集训文化，在解答中融入经验分享，促进成员交流</p>
      <div class="search-box">
        <el-input v-model="searchText" placeholder="输入关键词快速查找问题..." clearable class="search-input">
          <template #prefix>
            <el-icon>
              <Search />
            </el-icon>
          </template>
        </el-input>
      </div>
    </div>
    <transition-group name="fade-slide" tag="div" class="category-tabs">
      <div v-for="tab in tabs" :key="tab.value"
        :class="{ 'active-tab': activeTab === tab.value, 'hover-tab': hoveringTab === tab.value }"
        @mouseenter="hoveringTab = tab.value" @mouseleave="hoveringTab = null" @click="activeTab = tab.value">
        {{ tab.label }}
      </div>
    </transition-group>
    <transition name="fade-stagger" mode="out-in">
      <div v-if="filteredFaqs.length > 0" key="content">
        <div class="faq-list">
          <div v-for="(faq, index) in filteredFaqs" :key="faq.id">
            <div class="faq-question" @click="toggleAnswer(index)" @mouseenter="handleQuestionHover(index)"
              @mouseleave="handleQuestionLeave(index)">
              {{ faq.question }}
              <i :class="['fa', expandedAnswers[index] ? 'fa-minus' : 'fa-plus']"></i>
            </div>
            <transition name="answer-slide">
              <div class="faq-answer" v-if="expandedAnswers[index]" @mouseenter="handleAnswerHover(index)"
                @mouseleave="handleAnswerLeave(index)">
                {{ faq.answer }}
              </div>
            </transition>
            <hr class="divider">
          </div>
        </div>
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
import { ref, computed, onMounted } from 'vue';
import request from '@/utils/request';
import { Search } from '@element-plus/icons-vue';

export default {
  name: 'FaqPage',
  components: { Search },
  setup() {
    const faqList = ref([]);
    const searchText = ref('');
    const activeTab = ref('all');
    const hoveringTab = ref(null);
    const tabs = [
      { label: '全部问题', value: 'all' },
      { label: '入队相关', value: 'join' },
      { label: '训练相关', value: 'training' },
      { label: '比赛相关', value: 'competition' },
    ];
    const expandedAnswers = ref({});
    const questionHoverStates = ref({});
    const answerHoverStates = ref({});

    const filteredFaqs = computed(() => {
      if (activeTab.value === 'all') {
        return faqList.value.filter((item) =>
          item.question.includes(searchText.value) ||
          item.answer.includes(searchText.value)
        );
      } else {
        return faqList.value.filter((item) =>
          item.category === activeTab.value && (
            item.question.includes(searchText.value) ||
            item.answer.includes(searchText.value)
          )
        );
      }
    });

    const getFaqList = async () => {
      try {
        const res = await request.get('/faq/list');
        faqList.value = res.data.data || [];
      } catch (error) {
        console.error('获取FAQ失败', error);
      }
    };

    const resetSearch = () => {
      searchText.value = '';
    };

    const toggleAnswer = (index) => {
      expandedAnswers.value[index] = !expandedAnswers.value[index];
    };

    const handleQuestionHover = (index) => {
      questionHoverStates.value[index] = true;
    };

    const handleQuestionLeave = (index) => {
      questionHoverStates.value[index] = false;
    };

    const handleAnswerHover = (index) => {
      answerHoverStates.value[index] = true;
    };

    const handleAnswerLeave = (index) => {
      answerHoverStates.value[index] = false;
    };

    onMounted(() => {
      getFaqList();
    });

    return {
      faqList,
      searchText,
      activeTab,
      hoveringTab,
      tabs,
      filteredFaqs,
      resetSearch,
      expandedAnswers,
      toggleAnswer,
      questionHoverStates,
      answerHoverStates,
      handleQuestionHover,
      handleQuestionLeave,
      handleAnswerHover,
      handleAnswerLeave,
    };
  },
};
</script>

<style scoped>
body {
  font-family: 'Inter', sans-serif;
  line-height: 1.6;
}

.faq-page {
  padding: 1.5rem;
  max-width: 1000px;
  margin: 0 auto;
  min-height: 80vh;
}

.header {
  text-align: center;
  margin-bottom: 1.5rem;
  background: linear-gradient(135deg, #8ab7df 0%);
  padding: 3rem 1.5rem;
  border-radius: 1.5rem;
  box-shadow: 0 8px 20px rgba(79, 172, 254, 0.15);
  color: white;
  /* 取消顶部背景边框加深 */
  border: 1px solid rgba(0, 123, 255, 0.3);
}

.title {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 1rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle {
  font-size: clamp(1rem, 3vw, 1.2rem);
  opacity: 0.95;
  font-weight: 300;
  margin-bottom: 1.5rem;
}

.search-box {
  max-width: 450px;
  margin: 0 auto;
}

.search-input {
  border-radius: 5px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 50px;
  padding: 0.3rem 1.5rem;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
  height: 35px;
}

.category-tabs {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 1.5rem;
  position: relative;
}

.category-tabs::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(to right, transparent, rgba(0, 123, 255, 0.3), transparent);
  z-index: 1;
}

.category-tabs div {
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  z-index: 2;
  border: 1px solid transparent;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.hover-tab {
  background-color: rgba(0, 123, 255, 0.05);
  border-color: rgba(0, 123, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1), 0 0 8px rgba(255, 255, 255, 0.8);
}

.active-tab {
  background-color: rgba(0, 123, 255, 0.1);
  border-color: rgba(0, 123, 255, 0.5);
  transform: translateY(-4px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15), 0 0 12px rgba(255, 255, 255, 0.8);
}

.faq-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.faq-question {
  font-size: 1rem;
  font-weight: 600;
  color: #333;
  background-color: white;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  padding: 0.5rem;
  border-image: linear-gradient(to right, #e0e0e0, #f0f0f0) 1;
  transition: all 0.3s ease;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
}

.faq-question:hover,
.faq-question[data-hover="true"] {
  background-color: #f5f5f5;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  border-color: #ccc;
  transform: scale(1.02);
}

.faq-answer {
  font-size: 0.9rem;
  color: #666;
  line-height: 1.6;
  background-color: white;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  padding: 0.5rem;
  border-image: linear-gradient(to right, #e0e0e0, #f0f0f0) 1;
  transition: all 0.3s ease;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.faq-answer:hover,
.faq-answer[data-hover="true"] {
  background-color: #f5f5f5;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  border-color: #ccc;
  transform: scale(1.02);
}

.divider {
  border: none;
  border-top: 1px solid #e0e0e0;
  margin: 1rem 0;
}

.empty-state {
  text-align: center;
  padding: 3rem 0;
  background: white;
  border-radius: 1.5rem;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.05);
}

/* 分类切换动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 问答板块动画 */
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

/* 答案折叠展开动画 */
.answer-slide-enter-active,
.answer-slide-leave-active {
  transition: max-height 0.3s ease-in-out, opacity 0.3s ease-in-out, transform 0.3s ease;
  max-height: 1000px;
  opacity: 1;
  transform: scaleY(1);
  transform-origin: top;
}

.answer-slide-enter-from,
.answer-slide-leave-to {
  max-height: 0;
  opacity: 0;
  transform: scaleY(0);
}

.list-move {
  transition: transform 0.3s ease;
}

@media (max-width: 768px) {
  .title {
    font-size: 2rem;
  }

  .faq-question {
    font-size: 0.9rem;
  }

  .faq-answer {
    font-size: 0.8rem;
  }
}
</style>