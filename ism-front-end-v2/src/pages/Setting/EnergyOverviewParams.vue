<template>
  <a-spin :spinning="loading">
    <a-form layout="vertical" class="energy-overview-form">
      <a-alert
        type="info"
        show-icon
        :message="$t('SystemParams.EnergyOverview.hint')"
        class="energy-overview-form__hint"
      />
      <a-form-item
        v-for="field in keywordFields"
        :key="field.key"
        :label="$t(field.label)"
        v-bind="formLayout"
      >
        <a-input
          v-model="form[field.key]"
          :placeholder="$t('SystemParams.EnergyOverview.keywordPlaceholder')"
        />
      </a-form-item>
      <a-form-item :label="$t('SystemParams.EnergyOverview.bucketMinutes')" v-bind="formLayout">
        <a-input-number v-model="form.bucketMinutes" :min="1" :max="5" :precision="0" />
      </a-form-item>
      <a-form-item :label="$t('SystemParams.EnergyOverview.sampleIntervalSeconds')" v-bind="formLayout">
        <a-input-number v-model="form.sampleIntervalSeconds" :min="60" :max="300" :step="10" :precision="0" />
      </a-form-item>
      <a-form-item :wrapper-col="{ span: 10, offset: 3 }">
        <a-button type="primary" :loading="saving" @click="save">
          {{ $t('AlarmTips.Save') }}
        </a-button>
        <a-button class="energy-overview-form__scan" :loading="scanning" @click="scan">
          {{ $t('SystemParams.EnergyOverview.rescan') }}
        </a-button>
      </a-form-item>
      <div class="energy-overview-coverage">
        <div class="energy-overview-coverage__title">{{ $t('SystemParams.EnergyOverview.coverage') }}</div>
        <div class="energy-overview-coverage__metrics">
          <span>{{ $t('SystemParams.EnergyOverview.totalDevices') }}：{{ coverage.totalDevices }}</span>
          <span>{{ $t('SystemParams.EnergyOverview.eligibleDevices') }}：{{ coverage.eligibleDevices }}</span>
          <span>{{ $t('SystemParams.EnergyOverview.missingDevices') }}：{{ coverage.missingDevices }}</span>
          <span>{{ $t('SystemParams.EnergyOverview.ambiguousDevices') }}：{{ coverage.ambiguousDevices }}</span>
        </div>
        <div v-if="coverage.missingExamples.length" class="energy-overview-coverage__examples">
          {{ $t('SystemParams.EnergyOverview.missingExamples') }}：{{ coverage.missingExamples.join('、') }}
        </div>
        <div v-if="coverage.ambiguousExamples.length" class="energy-overview-coverage__examples">
          {{ $t('SystemParams.EnergyOverview.ambiguousExamples') }}：{{ coverage.ambiguousExamples.join('、') }}
        </div>
      </div>
    </a-form>
  </a-spin>
</template>

<script>
import {
  getEnergyOverviewConfig,
  getEnergyOverviewCandidates,
  refreshEnergyOverviewStats,
  saveEnergyOverviewConfig,
} from '@/services/energyOverview'

const EMPTY_FORM = {
  activeKeywords: '总有功功率|有功功率',
  reactiveKeywords: '总无功功率|无功功率',
  apparentKeywords: '总视在功率|视在功率',
  energyKeywords: '正有功电度|正向有功电能',
  bucketMinutes: 5,
  sampleIntervalSeconds: 60,
}

function resultOf(response) {
  return (response && response.data && response.data.result) || {}
}

export default {
  name: 'EnergyOverviewParams',
  i18n: require('../../i18n/language'),
  data() {
    return {
      loading: false,
      saving: false,
      scanning: false,
      form: { ...EMPTY_FORM },
      coverage: {
        totalDevices: 0,
        eligibleDevices: 0,
        missingDevices: 0,
        ambiguousDevices: 0,
        missingExamples: [],
        ambiguousExamples: [],
      },
      formLayout: {
        labelCol: { span: 3 },
        wrapperCol: { span: 10 },
      },
      keywordFields: [
        { key: 'activeKeywords', label: 'SystemParams.EnergyOverview.activeKeywords' },
        { key: 'reactiveKeywords', label: 'SystemParams.EnergyOverview.reactiveKeywords' },
        { key: 'apparentKeywords', label: 'SystemParams.EnergyOverview.apparentKeywords' },
        { key: 'energyKeywords', label: 'SystemParams.EnergyOverview.energyKeywords' },
      ],
    }
  },
  created() {
    this.load()
  },
  methods: {
    load() {
      this.loading = true
      Promise.all([getEnergyOverviewCandidates(), getEnergyOverviewConfig()])
        .then(([candidateResponse, configResponse]) => {
          this.applyCoverage(resultOf(candidateResponse))
          const config = resultOf(configResponse)
          this.form = {
            activeKeywords: config.activeKeywords || EMPTY_FORM.activeKeywords,
            reactiveKeywords: config.reactiveKeywords || EMPTY_FORM.reactiveKeywords,
            apparentKeywords: config.apparentKeywords || EMPTY_FORM.apparentKeywords,
            energyKeywords: config.energyKeywords || EMPTY_FORM.energyKeywords,
            bucketMinutes: Number(config.bucketMinutes || config.BucketMinutes || 5),
            sampleIntervalSeconds: Number(config.sampleIntervalSeconds || config.SampleIntervalSeconds || 60),
          }
        })
        .catch(() => this.$message.error(this.$t('SystemParams.EnergyOverview.loadFailed')))
        .finally(() => {
          this.loading = false
        })
    },
    applyCoverage(value) {
      const coverage = value || {}
      this.coverage = {
        totalDevices: Number(coverage.totalDevices || 0),
        eligibleDevices: Number(coverage.eligibleDevices || 0),
        missingDevices: Number(coverage.missingDevices || 0),
        ambiguousDevices: Number(coverage.ambiguousDevices || 0),
        missingExamples: coverage.missingExamples || [],
        ambiguousExamples: coverage.ambiguousExamples || [],
      }
    },
    scan() {
      this.scanning = true
      return getEnergyOverviewCandidates()
        .then(response => this.applyCoverage(resultOf(response)))
        .catch(() => this.$message.error(this.$t('SystemParams.EnergyOverview.loadFailed')))
        .finally(() => {
          this.scanning = false
        })
    },
    validate() {
      if (this.keywordFields.some(field => !String(this.form[field.key] || '').trim())) {
        this.$message.warning(this.$t('SystemParams.EnergyOverview.required'))
        return false
      }
      if (this.form.bucketMinutes < 1 || this.form.bucketMinutes > 5
        || this.form.sampleIntervalSeconds < 60 || this.form.sampleIntervalSeconds > 300) {
        this.$message.warning(this.$t('SystemParams.EnergyOverview.rangeError'))
        return false
      }
      return true
    },
    save() {
      if (!this.validate()) return
      this.saving = true
      saveEnergyOverviewConfig({ ...this.form })
        .then(response => {
          if (response.data && response.data.code === 0) {
            this.$message.success(this.$t('SystemParams.EnergyOverview.saveSuccess'))
            return Promise.all([
              refreshEnergyOverviewStats().catch(() => {}),
              this.scan(),
            ])
          } else {
            this.$message.error(this.$t('SystemParams.EnergyOverview.saveFailed'))
          }
        })
        .catch(() => this.$message.error(this.$t('SystemParams.EnergyOverview.saveFailed')))
        .finally(() => {
          this.saving = false
        })
    },
  },
}
</script>

<style scoped>
.energy-overview-form {
  padding: 10px;
}
.energy-overview-form__hint {
  margin-bottom: 16px;
}
.energy-overview-form__scan {
  margin-left: 12px;
}
.energy-overview-coverage {
  margin: 8px 0 0 12.5%;
  width: 41.6667%;
  padding: 16px;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  background: #fafafa;
}
.energy-overview-coverage__title {
  margin-bottom: 10px;
  font-weight: 600;
}
.energy-overview-coverage__metrics {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px 16px;
}
.energy-overview-coverage__examples {
  margin-top: 8px;
  color: rgba(0, 0, 0, 0.65);
  word-break: break-all;
}
</style>
