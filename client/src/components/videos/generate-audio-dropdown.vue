<script setup lang="ts">
import type { FormInst, FormItemRule, SelectGroupOption } from 'naive-ui'
import { reactive } from 'vue'
import { ref } from 'vue'

import { InformationCircleOutline as InfoIcon } from '@vicons/ionicons5'

import { PLAYHT_DEFAULT_SPEED, PLAYHT_DEFAULT_VOICE } from '@/config/constants'
import { PlayhtVoices } from '@/data/playht'

const formRef = ref<FormInst>()

const emit = defineEmits<{
  (name: 'submit', voice: string, speed: number, seed: number): void
}>()

const rules = {
  voice: {
    required: true,
    trigger: ['blur', 'input'],
    message: 'Please select voice',
  },
  speed: {
    required: true,
    trigger: ['blur', 'input'],
    validator(rule: FormItemRule, value: number) {
      if (value < 0.1) {
        return new Error('Speed is too low')
      }
      if (value > 5) {
        return new Error('Speed must be less than or equal to 5')
      }
      if (!/^[\d]+\.?[\d]{0,2}$/.test(String(value))) {
        return new Error('Please enter a valid number')
      }
      return true
    },
  },
  seed: {
    trigger: ['blur', 'input'],
    validator(rule: FormItemRule, value: number) {
      if (!/^\d{0,10}$/.test(String(value))) {
        return new Error('Invalid value')
      }
      return true
    },
  },
}

// Prepare options for the select grouped by accent
const options = PlayhtVoices.reduce<SelectGroupOption[]>((acc, voice) => {
  const group = acc.find((item) => item.label === voice.accent)

  if (group) {
    group.children!.push({
      label: voice.label,
      value: voice.id,
    })
  } else {
    acc.push({
      type: 'group',
      label: voice.accent,
      key: voice.accent,
      children: [
        {
          label: voice.label,
          value: voice.id,
        },
      ],
    })
  }

  return acc
}, [])

const model = reactive({
  voice: PLAYHT_DEFAULT_VOICE,
  speed: PLAYHT_DEFAULT_SPEED,
  seedType: 'random' as 'random' | 'custom',
  seed: 0,
})

const seedTypes = [
  { label: 'Random', value: 'random' },
  { label: 'Custom', value: 'custom' },
]

const methods = {
  async submit() {
    if (!formRef.value) return
    try {
      await formRef.value.validate()
      methods.generate()
    } catch {
      // noop
    }
  },
  generate() {
    let seed = model.seed ?? 0
    if (model.seedType === 'random') {
      seed = Math.floor(Math.random() * 1000000)
    }
    emit('submit', model.voice, model.speed, seed)
  },
}
</script>

<template>
  <div class="generate-audio-dropdown">
    <n-form
      ref="formRef"
      :model="model"
      :rules="rules"
      label-placement="top"
      size="small"
      :label-width="120"
      style="min-width: 200px"
    >
      <n-form-item label="Voice" path="voice">
        <n-select v-model:value="model.voice" placeholder="Select" :options="options" filterable />
      </n-form-item>
      <n-form-item label="Speed" path="speed">
        <n-input-number v-model:value="model.speed" :step="0.01" style="max-width: 95px" />
      </n-form-item>
      <n-form-item path="seed">
        <template #label>
          Seed
          <n-tooltip trigger="hover" placement="bottom">
            <template #trigger>
              <InfoIcon class="info-icon" />
            </template>
            Used to control the reproducibility of the generated audio. A fixed seed should always
            generate the exact same audio file.
          </n-tooltip>
        </template>
        <div>
          <n-radio-group v-model:value="model.seedType" name="radiogroup">
            <n-space>
              <n-radio
                v-for="{ value, label } in seedTypes"
                :key="value"
                :value="value"
                :label="label"
              />
            </n-space>
          </n-radio-group>

          <p v-if="model.seedType === 'custom'">
            <n-input-number v-model:value="model.seed" clearable placeholder="12345" />
          </p>

          <p v-if="model.seedType === 'random'" class="description">
            Random seed will produce a different audio file every time.
          </p>
        </div>
      </n-form-item>
    </n-form>

    <div style="display: flex; justify-content: flex-end">
      <n-button type="success" size="small" @click="methods.submit">Generate</n-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.generate-audio-dropdown {
  padding: 20px 24px;
  width: 280px;
}

.info-icon {
  width: 16px;
  vertical-align: middle;
  position: relative;
  top: -1px;
}

.description {
  font-size: 90%;
}
</style>
