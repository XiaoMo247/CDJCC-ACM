<template>
  <div class="honor-wall page-honor">
    <div class="page-header page-hero">
      <h1 class="page-title">
        <span class="page-title-icon"><i class="fas fa-star"></i></span>
        荣誉墙
      </h1>
      <p class="page-subtitle">记录每一次突破与成长，向优秀的队员致敬</p>
    </div>

    <div class="filters">
      <!-- 赛事筛选（按钮组） -->
      <div class="competition-buttons">
        <el-button 
          v-for="comp in competitionOptions" 
          :key="comp.value"
          :type="selectedCompetition === comp.value ? 'primary' : 'default'"
          @click="selectedCompetition = comp.value === '全部' ? '' : comp.value"
        >
          {{ comp.label }}
        </el-button>
      </div>

      <!-- 其他筛选（下拉框+姓名搜索） -->
      <div class="dropdown-filters">
        <el-select v-model="searchYear" placeholder="全部年份" class="filter-item" clearable>
          <el-option v-for="year in yearOptions" :key="year" :label="year" :value="year" />
        </el-select>

        <el-select v-model="selectedLevel" placeholder="全部级别" class="filter-item" clearable>
          <el-option v-for="level in levelOptions" :key="level" :label="level" :value="level" />
        </el-select>

        <el-select v-model="searchGrade" placeholder="全部年级" class="filter-item" clearable>
          <el-option v-for="grade in gradeOptions" :key="grade.value" :label="grade.label" :value="grade.value" />
        </el-select>

        <el-select v-model="searchMajor" placeholder="全部专业" class="filter-item" clearable>
          <el-option v-for="major in majorOptions" :key="major" :label="major" :value="major" />
        </el-select>
      </div>
    </div>

    <div class="honor-list">
      <h2 class="list-title">荣誉列表</h2>
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="6" animated />
      </div>
      <div v-else-if="filteredHonors.length === 0" class="no-data">
        <el-empty description="暂无荣誉数据" />
      </div>
      <div v-else class="honor-grid">
        <div v-for="honor in filteredHonors" :key="honor.id" class="honor-item">
          <img :src="defaultAvatar" :alt="honor.name" class="honor-avatar" />
          <div class="info">
            <h3>{{ honor.name }}</h3>
            <p>年份: {{ honor.year }}</p>
            <p>赛事: {{ honor.contest }}</p>
            <p>级别: {{ honor.level }}</p>
            <p>年级: {{ getGradeText(honor.grade) }}</p>
            <p>专业: {{ honor.major }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import request from '@/utils/request';

export default {
  name: 'HonorWall',
  data() {
    return {
      honors: [],
      loading: false,
      selectedCompetition: '',
      searchYear: '',
      selectedLevel: '',
      searchGrade: '',
      searchMajor: '',
      defaultAvatar: import.meta.env.VITE_DEFAULT_AVATAR || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
      
      // 筛选选项
      competitionOptions: [
        { label: '全部', value: '' },
        { label: 'ICPC', value: 'ICPC' },
        { label: 'CCPC', value: 'CCPC' },
        { label: '蓝桥杯', value: '蓝桥杯' },
        { label: '百度之星', value: '百度之星' }
      ],
      levelOptions: ['金奖', '银奖', '铜奖', '国一', '国二', '国三', '省一', '省二', '省三'],
      gradeOptions: [
        { label: '大一', value: 1 },
        { label: '大二', value: 2 },
        { label: '大三', value: 3 },
        { label: '大四', value: 4 }
      ],
    };
  },
  computed: {
    // 从数据中提取唯一值作为筛选选项
    yearOptions() {
      return [...new Set(this.honors.map(h => h.year))].sort((a, b) => b - a);
    },
    majorOptions() {
      return [...new Set(this.honors.map(h => h.major))].sort();
    },
    
    // 过滤后的荣誉数据
    filteredHonors() {
      return this.honors.filter(honor => {
        return (
          (this.selectedCompetition === '' || honor.contest === this.selectedCompetition) &&
          (this.searchYear === '' || honor.year.toString() === this.searchYear.toString()) &&
          (this.selectedLevel === '' || honor.level === this.selectedLevel) &&
          (this.searchGrade === '' || honor.grade.toString() === this.searchGrade.toString()) &&
          (this.searchMajor === '' || honor.major === this.searchMajor)
        );
      });
    }
  },
  mounted() {
    this.fetchHonors();
  },
  methods: {
    async fetchHonors() {
      this.loading = true;
      try {
        const res = await request.get('/honor');
        this.honors = res.data.data || [];
      } catch (err) {
        console.error('获取荣誉数据失败:', err);
        this.$message.error('获取荣誉数据失败');
      } finally {
        this.loading = false;
      }
    },
    
    getGradeText(grade) {
      const grades = ['', '大一', '大二', '大三', '大四'];
      return grades[grade] || '未知';
    },
    
    getInitials(name) {
      if (!name) return '';
      const names = name.split(' ');
      if (names.length > 1) {
        return names[0][0] + names[names.length - 1][0];
      }
      return name.length > 1 ? name.substring(0, 2) : name;
    }
  }
};
</script>

<style scoped>
.page-honor {
  padding: 0;
}

.page-hero {
  background: linear-gradient(135deg, #d8b4fe 0%, #c084fc 100%);
  padding: 60px 20px;
  text-align: center;
  color: #6b21a8;
  margin-bottom: 60px;
}

.page-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0 0 15px 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
}

.page-title-icon {
  font-size: 3.5rem;
}

.page-subtitle {
  font-size: 1.25rem;
  opacity: 0.95;
  margin: 0;
}

.honor-wall {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px 60px;
}

.filters {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 50px;
}

.competition-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 25px;
  padding-bottom: 25px;
  border-bottom: 1px solid #f0f0f0;
}

.competition-buttons .el-button {
  border-radius: 20px;
  padding: 8px 20px;
}

.dropdown-filters {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
  margin-bottom: 0;
}

.filter-item {
  min-width: 140px;
}

.honor-list {
  margin-top: 40px;
}

.list-title {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 40px;
  color: #333;
  display: flex;
  align-items: center;
  gap: 12px;
}

.list-title::before {
  content: '✨';
  font-size: 2rem;
}

.honor-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 32px;
  padding: 0;
}

.honor-item {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border-left: 3px solid #d8b4fe;
  position: relative;
  overflow: hidden;
}

.honor-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(135deg, rgba(216, 180, 254, 0.02) 0%, transparent 100%);
  pointer-events: none;
}

.honor-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(216, 180, 254, 0.1);
  border-left-color: #c084fc;
}

.honor-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #d8b4fe;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(216, 180, 254, 0.15);
}

.honor-item img {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #d8b4fe;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(216, 180, 254, 0.15);
}

.info {
  flex: 1;
  width: 100%;
  position: relative;
  z-index: 1;
}

.info h3 {
  margin: 0 0 12px 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d2d2d;
}

.info p {
  margin: 8px 0;
  font-size: 0.9rem;
  color: #666;
  line-height: 1.5;
}

.info p::before {
  content: attr(data-icon);
  margin-right: 6px;
  font-weight: 600;
  color: #a855f7;
}

.loading-container {
  padding: 40px 20px;
}

.no-data {
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  text-align: center;
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }

  .page-subtitle {
    font-size: 1rem;
  }

  .honor-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .dropdown-filters {
    flex-direction: column;
  }

  .filter-item {
    width: 100%;
  }
}
</style>
