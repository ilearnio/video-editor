<script setup lang="ts">
import { c, useNotification } from 'naive-ui'
import type { FormInst, FormRules, FormValidationError } from 'naive-ui'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { t } from '@/i18n'
import { authenticateUser } from '@/services/api/users'
import { useUserStore } from '@/stores/user'

interface ModelType {
  email: string | null
  password: string | null
}

const router = useRouter()
const notification = useNotification()
const formRef = ref<FormInst | null>(null)
const formLoading = ref(false)
const userStore = useUserStore()
const formValue = ref<ModelType>({
  email: null,
  password: null,
})
const rules: FormRules = {
  email: {
    required: true,
    type: 'email',
    message: t('login.email.required'),
    trigger: ['input'],
  },
  password: {
    required: true,
    message: t('login.password.required'),
    trigger: ['input'],
  },
}
const login = async (e: MouseEvent) => {
  e.preventDefault()
  formLoading.value = true
  formRef.value?.validate(async (errors: Array<FormValidationError> | undefined) => {
    if (!errors) {
      const { email, password } = formValue.value
      try {
        const { user } = await authenticateUser(email!, password!)
        notification.success({
          closable: true,
          duration: 5000,
          content: t('notification.content.auth'),
          meta: t('notification.login', { username: user.username }),
        })
        router.push('/')
      } catch (err) {
        notification.error({
          duration: 3000,
          content: t('notification.content.auth'),
          meta: (err as any).response?.data?.message ?? t('notification.unknownerror'),
        })
      } finally {
        formLoading.value = false
      }
    } else {
      formLoading.value = false
    }
  })
}
</script>

<template>
  <n-h1 style="--font-size: 60px; --font-weight: 100">
    <img alt="Devstack logo" class="logo" src="@/assets/logo.png" width="100" />
    Devstack
  </n-h1>
  <n-card size="large" style="--padding-bottom: 30px">
    <n-h2 style="--font-weight: 400">{{ $t('login.signin') }}</n-h2>
    <n-space vertical>
      <n-form ref="formRef" :model="formValue" :rules="rules">
        <n-form-item-row :label="$t('login.email')" path="email">
          <n-input
            v-model:value="formValue.email"
            placeholder=""
            :input-props="{ type: 'email', autocomplete: 'off' }"
          />
        </n-form-item-row>
        <n-form-item-row path="password" :label="$t('login.password')">
          <n-input
            v-model:value="formValue.password"
            placeholder=""
            type="password"
            @keydown.enter.prevent
          />
        </n-form-item-row>
        <n-button type="primary" size="large" block :loading="formLoading" @click="login">
          <template #icon>
            <Icon type="login" />
          </template>
          {{ $t('login.signin') }}
        </n-button>
      </n-form>
    </n-space>
  </n-card>
</template>
<style scoped lang="scss">
.n-h1 {
  margin: 20vh auto 20px;
  text-align: center;
  text-align: -webkit-center;
  letter-spacing: 5px;
  opacity: 0.8;

  img {
    display: block;
  }
}

.n-card {
  margin: 0 auto;
  max-width: 380px;
  box-shadow: var(--box-shadow);
}
</style>
