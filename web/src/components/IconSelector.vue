<template>
  <div class="icon-selector">
    <el-popover
      placement="bottom-start"
      :width="360"
      trigger="click"
      popper-class="icon-popover"
    >
      <template #reference>
        <el-input
          v-model="displayValue"
          placeholder="请选择图标"
          readonly
          @click="showPopover = true"
          class="cursor-pointer"
        >
          <template #prefix>
            <el-icon v-if="modelValue" class="text-lg">
              <component :is="modelValue" />
            </el-icon>
            <el-icon v-else class="text-gray-400">
              <Plus />
            </el-icon>
          </template>
          <template #suffix>
            <el-icon class="text-gray-400">
              <ArrowDown />
            </el-icon>
          </template>
        </el-input>
      </template>
      
      <div class="icon-panel">
        <div class="mb-3">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索图标..."
            size="small"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        
        <div class="icons-grid">
          <div
            v-for="icon in filteredIcons"
            :key="icon.name"
            class="icon-item"
            :class="{ 'selected': modelValue === icon.name }"
            @click="selectIcon(icon.name)"
          >
            <el-icon class="icon-display">
              <component :is="icon.component" />
            </el-icon>
            <span class="icon-name">{{ icon.name }}</span>
          </div>
        </div>
        
        <div class="mt-3 text-center" v-if="modelValue">
          <el-button size="small" @click="clearIcon()">清除选择</el-button>
        </div>
      </div>
    </el-popover>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
// 根据Element Plus官方文档，只导入确实存在的图标
import { 
  Plus, 
  ArrowDown, 
  Search,
  House,
  User,
  Setting,
  Menu,
  UserFilled,
  Document,
  Folder,
  FolderOpened,
  Files,
  Edit,
  Delete,
  View,
  Share,
  Download,
  Upload,
  Refresh,
  Close,
  Check,
  Star,
  Location,
  Phone,
  Message,
  ChatDotRound,
  Bell,
  Warning,
  QuestionFilled,
  InfoFilled,
  SuccessFilled,
  CircleClose,
  Lock,
  Unlock,
  Key,
  Tools,
  DataLine,
  PieChart,
  Histogram,
  TrendCharts,
  Monitor,
  Camera,
  Picture,
  VideoCamera,
  Calendar,
  Clock,
  Flag,
  Trophy,
  Present,
  ShoppingCart,
  Goods,
  Money,
  CreditCard,
  Van,
  Guide,
  Connection,
  Link,
  Management,
  Operation,
  Promotion,
  Rank,
  CirclePlus,
  CircleCheck,
  Loading,
  Finished,
  Position,
  Service,
  Reading,
  Notebook,
  School,
  OfficeBuilding,
  Suitcase,
  HomeFilled,
  More,
  Compass,
  Filter,
  Switch,
  Select,
  CloseBold,
  EditPen,
  MessageBox,
  TurnOff,
  Crop,
  SwitchButton,
  Open,
  Remove,
  ZoomOut,
  ZoomIn,
  CircleCheckFilled,
  WarningFilled,
  CircleCloseFilled,
  HelpFilled,
  StarFilled,
  Comment,
  Grid,
  DeleteFilled,
  RemoveFilled,
  CirclePlusFilled
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const searchKeyword = ref('')

// 根据Element Plus官方文档整理的图标列表，确保所有图标都是存在的
const commonIcons = [
  // 系统类图标
  { name: 'House', component: House },
  { name: 'HomeFilled', component: HomeFilled },
  { name: 'User', component: User },
  { name: 'UserFilled', component: UserFilled },
  { name: 'Setting', component: Setting },
  { name: 'Menu', component: Menu },
  { name: 'More', component: More },
  { name: 'Search', component: Search },
  { name: 'Plus', component: Plus },
  { name: 'Check', component: Check },
  { name: 'Close', component: Close },
  { name: 'CloseBold', component: CloseBold },
  { name: 'Loading', component: Loading },
  { name: 'Finished', component: Finished },
  { name: 'Refresh', component: Refresh },
  
  // 文档类图标
  { name: 'Document', component: Document },
  { name: 'Folder', component: Folder },
  { name: 'FolderOpened', component: FolderOpened },
  { name: 'Files', component: Files },
  { name: 'Notebook', component: Notebook },
  { name: 'Reading', component: Reading },
  
  // 操作类图标
  { name: 'Edit', component: Edit },
  { name: 'EditPen', component: EditPen },
  { name: 'Delete', component: Delete },
  { name: 'DeleteFilled', component: DeleteFilled },
  { name: 'View', component: View },
  { name: 'Share', component: Share },
  { name: 'Download', component: Download },
  { name: 'Upload', component: Upload },
  { name: 'Remove', component: Remove },
  { name: 'RemoveFilled', component: RemoveFilled },
  { name: 'ZoomIn', component: ZoomIn },
  { name: 'ZoomOut', component: ZoomOut },
  { name: 'Crop', component: Crop },
  { name: 'Open', component: Open },
  
  // 状态类图标
  { name: 'CirclePlus', component: CirclePlus },
  { name: 'CirclePlusFilled', component: CirclePlusFilled },
  { name: 'CircleCheck', component: CircleCheck },
  { name: 'CircleCheckFilled', component: CircleCheckFilled },
  { name: 'CircleClose', component: CircleClose },
  { name: 'CircleCloseFilled', component: CircleCloseFilled },
  { name: 'Warning', component: Warning },
  { name: 'WarningFilled', component: WarningFilled },
  { name: 'QuestionFilled', component: QuestionFilled },
  { name: 'InfoFilled', component: InfoFilled },
  { name: 'SuccessFilled', component: SuccessFilled },
  { name: 'HelpFilled', component: HelpFilled },
  { name: 'Star', component: Star },
  { name: 'StarFilled', component: StarFilled },
  { name: 'Flag', component: Flag },
  
  // 通讯类图标
  { name: 'Message', component: Message },
  { name: 'MessageBox', component: MessageBox },
  { name: 'ChatDotRound', component: ChatDotRound },
  { name: 'Comment', component: Comment },
  { name: 'Phone', component: Phone },
  { name: 'Bell', component: Bell },
  
  // 位置类图标
  { name: 'Location', component: Location },
  { name: 'Position', component: Position },
  { name: 'Compass', component: Compass },
  { name: 'Guide', component: Guide },
  
  // 安全类图标
  { name: 'Lock', component: Lock },
  { name: 'Unlock', component: Unlock },
  { name: 'Key', component: Key },
  
  // 工具类图标
  { name: 'Tools', component: Tools },
  { name: 'Filter', component: Filter },
  { name: 'Switch', component: Switch },
  { name: 'SwitchButton', component: SwitchButton },
  { name: 'Select', component: Select },
  { name: 'TurnOff', component: TurnOff },
  
  // 商务类图标
  { name: 'Management', component: Management },
  { name: 'Operation', component: Operation },
  { name: 'Service', component: Service },
  { name: 'Promotion', component: Promotion },
  { name: 'Rank', component: Rank },
  { name: 'Connection', component: Connection },
  { name: 'Link', component: Link },
  
  // 数据类图标
  { name: 'DataLine', component: DataLine },
  { name: 'PieChart', component: PieChart },
  { name: 'Histogram', component: Histogram },
  { name: 'TrendCharts', component: TrendCharts },
  { name: 'Grid', component: Grid },
  
  // 时间类图标
  { name: 'Calendar', component: Calendar },
  { name: 'Clock', component: Clock },
  
  // 媒体类图标
  { name: 'Monitor', component: Monitor },
  { name: 'Camera', component: Camera },
  { name: 'Picture', component: Picture },
  { name: 'VideoCamera', component: VideoCamera },
  
  // 商品类图标
  { name: 'ShoppingCart', component: ShoppingCart },
  { name: 'Goods', component: Goods },
  { name: 'Money', component: Money },
  { name: 'CreditCard', component: CreditCard },
  { name: 'Present', component: Present },
  { name: 'Trophy', component: Trophy },
  
  // 建筑类图标
  { name: 'School', component: School },
  { name: 'OfficeBuilding', component: OfficeBuilding },
  { name: 'Van', component: Van },
  
  // 物品类图标
  { name: 'Suitcase', component: Suitcase }
]

const displayValue = computed(() => {
  return props.modelValue || ''
})

const filteredIcons = computed(() => {
  if (!searchKeyword.value) {
    return commonIcons
  }
  return commonIcons.filter(icon => 
    icon.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const selectIcon = (iconName) => {
  emit('update:modelValue', iconName)
}

const clearIcon = () => {
  emit('update:modelValue', '')
}
</script>

<style lang="scss" scoped>
.icon-selector {
  .cursor-pointer {
    cursor: pointer;
  }
}

.icon-panel {
  .icons-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 8px;
    max-height: 300px;
    overflow-y: auto;
    
    .icon-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 8px 4px;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s;
      min-height: 60px;
      
      &:hover {
        border-color: #409eff;
        background-color: #f0f8ff;
      }
      
      &.selected {
        border-color: #409eff;
        background-color: #e6f7ff;
        box-shadow: 0 0 0 1px #409eff;
      }
      
      .icon-display {
        font-size: 20px;
        margin-bottom: 4px;
        color: #606266;
      }
      
      .icon-name {
        font-size: 10px;
        color: #909399;
        text-align: center;
        line-height: 1.2;
        word-break: break-all;
      }
    }
  }
  
  .el-input {
    margin-bottom: 10px;
  }
}

.icon-popover {
  .icon-panel {
    padding: 8px;
  }
}
</style> 