<script setup lang="ts">
import debounce from 'just-debounce'
import { type DropdownOption, NIcon, NTag } from 'naive-ui'
import type {
  FilterState,
  SortState,
  TableBaseColumn,
  TableColumns,
} from 'naive-ui/es/data-table/src/interface'
import type { RecordListQueryParams } from 'pocketbase'
import { type Component, computed, createVNode, h, onMounted, reactive } from 'vue'
import { RouterLink } from 'vue-router'

import { MenuOutline as MenuIcon, RemoveCircleOutline } from '@vicons/ionicons5'

import { formatTimestamp } from '@/helpers/date'
import { upperFirst } from '@/helpers/strings'
import type { Video } from '@/models/videos'
import { useVideoStore } from '@/stores/video'

import TableRowActions from './partials/table-row-actions.vue'

const videosStore = useVideoStore()

const renderIcon = (icon: Component) => () => {
  return h(NIcon, null, {
    default: () => h(icon),
  })
}

const columns: TableColumns<Video> = [
  {
    type: 'selection',
  },
  {
    title: 'Title',
    key: 'title',
    sorter: 'default',
    minWidth: 200,
    render(row) {
      return h(
        RouterLink,
        { to: { name: 'videos.edit', params: { id: row.id } }, class: 'title-link' },
        { default: () => row.title },
      )
    },
  },
  {
    title: 'Type',
    key: 'type',
    sorter: 'default',
    width: 120,
    render: (row) => upperFirst(row.type),
  },
  {
    title: 'Status',
    key: 'status',
    sorter: 'default',
    width: 120,
    filter: true,
    filterOptions: [
      {
        label: 'Draft',
        value: 'draft',
      },
      {
        label: 'Completed',
        value: 'completed',
      },
    ],
    render(row) {
      if (!row.status) return
      return h(
        NTag,
        { size: 'small', type: row.status === 'completed' ? 'success' : 'default' },
        { default: () => upperFirst(row.status) },
      )
    },
  },
  {
    title: 'Creation Date',
    key: 'created',
    sorter: 'default',
    width: 170,
    render: (row) => formatTimestamp(row.created!),
  },
  {
    title: 'Updated Date',
    key: 'updated',
    sorter: 'default',
    width: 170,
    render: (row) => formatTimestamp(row.updated!),
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 110,
    render(row) {
      return createVNode(TableRowActions, {
        video: row,
        onDelete: methods.fetchItems,
      })
    },
  },
]

const data = reactive({
  isLoading: false,
  searchValue: '',
  checkedRowKeys: [] as string[],
  page: 1,
  sortByColumn: '',
  sortOrder: '' as 'ascend' | 'descend' | '',
  filters: {} as Record<string, string[]>,
})

const getters = {
  pagination: computed(() => {
    return {
      pageSize: videosStore.state.list.perPage,
      pageCount: videosStore.state.list.totalPages,
      itemCount: videosStore.state.list.totalItems,
    }
  }),
  checkedRowActions: computed<DropdownOption[]>(() => [
    {
      label: 'Delete',
      key: 'delete',
      icon: renderIcon(RemoveCircleOutline),
    },
  ]),
}

const methods = {
  async fetchItems() {
    data.isLoading = true
    try {
      const params: RecordListQueryParams = {}

      if (data.searchValue) {
        params.filter = `title ?~ '${data.searchValue.replace("'", "\\'")}'`
      }

      if (data.page) {
        params.page = data.page
      }

      if (data.sortByColumn) {
        params.sort = data.sortByColumn
        if (data.sortOrder === 'descend') {
          params.sort = `-${params.sort}`
        }
      }

      Object.entries(data.filters).forEach(([key, values]) => {
        if (values.length) {
          params.filter = `${params.filter ? `${params.filter} && ` : ''}(${values
            .map((value) => `${key} = '${value.replace("'", "\\'")}'`)
            .join(' && ')})`
        }
      })

      await videosStore.actions.fetchItems(params)
    } finally {
      data.isLoading = false
    }
  },
  search: debounce(async () => {
    await methods.fetchItems()
  }, 400),
  handleSelectCheckedRowsAction(optionValue: string) {
    if (optionValue === 'delete') {
      videosStore.actions.deleteItems(data.checkedRowKeys)
    }
  },
  handleCheck(keys: string[]) {
    data.checkedRowKeys = keys
  },
  async handlePageChange(page: number) {
    data.page = page
    await methods.fetchItems()
  },
  async handleSortChange(sortState: SortState) {
    if (sortState.order === false) {
      data.sortByColumn = ''
      data.sortOrder = ''
    } else {
      data.sortByColumn = sortState.columnKey as string
      data.sortOrder = sortState.order
    }
    await methods.fetchItems()
  },
  async handleFiltersChange(filterState: FilterState, sourceColumn: TableBaseColumn) {
    const values = filterState[sourceColumn.key] as string[]
    data.filters[sourceColumn.key] = values
    await methods.fetchItems()
  },
}

onMounted(async () => {
  await methods.fetchItems()
})
</script>

<template>
  <n-space vertical :size="12">
    <n-space>
      <n-input v-model:value="data.searchValue" placeholder="Search" @input="methods.search" />
      <n-dropdown
        :options="getters.checkedRowActions.value"
        placement="bottom-start"
        @select="methods.handleSelectCheckedRowsAction"
      >
        <n-button :disabled="data.checkedRowKeys.length === 0">
          <menu-icon width="20" /> &nbsp; Actions
        </n-button>
      </n-dropdown>
    </n-space>

    <n-data-table
      class="data-table"
      remote
      :columns="columns"
      :data="videosStore.state.list.items"
      :pagination="getters.pagination.value"
      :loading="data.isLoading"
      :row-key="(row) => row.id"
      @update:checked-row-keys="methods.handleCheck"
      @update:page="methods.handlePageChange"
      @update:sorter="methods.handleSortChange"
      @update:filters="methods.handleFiltersChange"
    />
  </n-space>
</template>

<style scoped lang="scss">
.data-table :deep(.title-link) {
  color: inherit;
}
</style>
