import { createI18n } from 'vue-i18n'

import { useStorage } from '@vueuse/core'

import enUs from '@/locales/en-US.json'

const defaultLanguage = 'en-US'
export const languages = ['en-US']

// browser language
const language = useStorage('language', defaultLanguage)

// default language
const i18n = createI18n({
  locale: language.value,
  legacy: false,
  globalInjection: true,
  fallbackLocale: 'en-US',
  messages: {
    'en-US': enUs,
  },
})
export const { t } = i18n.global

export default i18n
