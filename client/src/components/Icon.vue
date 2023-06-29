<template>
  <n-icon v-bind="props"><icon /></n-icon>
</template>

<script lang="ts" setup>
import { NIcon } from 'naive-ui'

import {
  AddOutline as add,
  ChatbubbleEllipsesOutline as comments,
  AddCircle as create,
  SpeedometerOutline as dashboard,
  AlertCircleOutline as fallback,
  FlaskOutline as flask,
  LogoGithub as github,
  HelpCircleOutline as help,
  LogInOutline as login,
  LogOutOutline as logout,
  ImagesOutline as media,
  MoonOutline as moon,
  NotificationsOutline as notifications,
  NewspaperOutline as pages,
  HardwareChipOutline as plugins,
  CreateOutline as posts,
  RefreshOutline as refresh,
  SaveOutline as save,
  SettingsOutline as settings,
  SunnyOutline as sunny,
  ColorPaletteOutline as themes,
  BuildOutline as tools,
  PeopleOutline as users,
} from '@vicons/ionicons5'
import { Trash as destroy, Pencil as edit } from '@vicons/tabler'

const icons = {
  destroy,
  edit,
  moon,
  sunny,
  create,
  save,
  login,
  logout,
  refresh,
  dashboard,
  posts,
  pages,
  media,
  users,
  comments,
  themes,
  plugins,
  tools,
  settings,
  flask,
  add,
  help,
  notifications,
  github,
  fallback,
}

type IconTypes = keyof typeof icons

const props = defineProps({
  ...NIcon.props,
  type: {
    type: String,
    default: 'fallback',
  },
})

const icon = icons[props.type as IconTypes]
</script>
