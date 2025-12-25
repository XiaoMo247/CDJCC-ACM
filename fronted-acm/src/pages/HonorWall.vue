<template>
  <div class="honor-wall">
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

      <!-- 其他筛选（下拉框） -->
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

      <!-- 姓名搜索 -->
      <el-input v-model="searchName" placeholder="学生姓名" class="name-search" clearable />
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
      searchName: '',
      defaultAvatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
      
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
          (this.searchMajor === '' || honor.major === this.searchMajor) &&
          (this.searchName === '' || honor.name.includes(this.searchName))
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
.honor-wall {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.filters {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 30px;
}

.competition-buttons {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.dropdown-filters {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.filter-item {
  width: 150px;
}

.name-search {
  width: 200px;
  margin: 0 auto;
}

.honor-list {
  text-align: center;
}

.list-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #333;
}

.honor-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  justify-items: center;
}

.honor-item {
  border: 1px solid #ddd;
  padding: 15px;
  border-radius: 8px;
  display: flex;
  gap: 15px;
  width: 100%;
  max-width: 350px;
  transition: all 0.3s;
  background-color: white;
}

.honor-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: #409EFF;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}

.info {
  flex: 1;
  text-align: left;
}

.info h3 {
  margin: 0 0 10px 0;
  font-size: 18px;
}

.info p {
  margin: 5px 0;
  font-size: 14px;
  color: #666;
}

.loading-container {
  padding: 20px;
}

.no-data {
  padding: 40px 0;
}

.honor-avatar {
  width: 80px; 
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #f0f0f0;
}

.honor-item img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #f0f0f0;
}
</style>
