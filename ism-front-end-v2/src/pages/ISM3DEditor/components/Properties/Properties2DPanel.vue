<template>
  <div class="right-panel">
    <div class="no-select-tip" v-if="!currentObject">
      <i class="fas fa-mouse-pointer"></i>
      {{ $t('ISM3DEditor.pleaseSelectObject') }}
    </div>

    <template v-else>
      <a-tabs v-model="activeKey" size="small" class="prop-tabs">
        <a-tab-pane key="style" :tab="$t('displayConfig.Properties.TabHeaterStyle')" :forceRender="true">
          <div class="tab-section">
            <div class="tab-section-title">{{ $t('displayConfig.Properties.TabHeaterStyle') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
              <a-form-item :label="$t('displayConfig.Properties.ComponentName')">
                <a-input :value="currentObject.name" @change="onRootPropChange('name', $event)" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentIsVisible')">
                <a-switch :checked="style.visible !== 0" @change="onStyleSwitch('visible', $event ? 1 : 0)" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentX')">
                <a-input-number :value="n(pos.x, 0)" @change="onPositionChange('x', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentY')">
                <a-input-number :value="n(pos.y, 0)" @change="onPositionChange('y', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentWith')">
                <a-input-number :value="n(pos.w, 160)" :min="24" @change="onPositionChange('w', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentHeight')">
                <a-input-number :value="n(pos.h, 80)" :min="24" @change="onPositionChange('h', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.Transform')">
                <a-input-number :value="n(style.transform, 0)" @change="onStyleChange('transform', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.opacity')">
                <a-input-number :value="n(style.opacity, 1)" :min="0" :max="1" :step="0.1" @change="onStyleChange('opacity', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.zIndex')">
                <a-input-number :value="n(style.zIndex, 0)" @change="onStyleChange('zIndex', $event)" style="width:100%" />
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section">
            <div class="tab-section-title">{{ $t('displayConfig.Properties.TabHeaterStyle') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
              <a-form-item v-if="style.backColor !== undefined" :label="$t('displayConfig.Properties.ComponentBackColor')">
                <input type="color" :value="normColor(style.backColor)" @input="onStyleChange('backColor', $event.target.value)" class="cip" />
              </a-form-item>
              <a-form-item v-if="style.foreColor !== undefined" :label="$t('displayConfig.Properties.ComponentForeColor')">
                <input type="color" :value="normColor(style.foreColor)" @input="onStyleChange('foreColor', $event.target.value)" class="cip" />
              </a-form-item>
              <a-form-item v-if="style.borderColor !== undefined" :label="$t('displayConfig.Properties.BorderColor')">
                <input type="color" :value="normColor(style.borderColor)" @input="onStyleChange('borderColor', $event.target.value)" class="cip" />
              </a-form-item>
              <a-form-item v-if="style.borderWidth !== undefined" :label="$t('displayConfig.Properties.BorderWidth')">
                <a-input-number :value="n(style.borderWidth, 0)" @change="onStyleChange('borderWidth', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item v-if="style.BorderEdges !== undefined" :label="$t('displayConfig.Properties.BorderEdges')">
                <a-input-number :value="n(style.BorderEdges, 0)" @change="onStyleChange('BorderEdges', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item v-if="style.text !== undefined" :label="$t('displayConfig.Properties.ComponentText')">
                <a-textarea :value="style.text" @change="onStyleInput('text', $event)" />
              </a-form-item>
              <a-form-item v-if="style.fontSize !== undefined" :label="$t('displayConfig.Properties.ComponentFontSize')">
                <a-input-number :value="n(style.fontSize, 14)" @change="onStyleChange('fontSize', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item v-if="style.fontWeight !== undefined" :label="$t('displayConfig.Properties.ComponentFontBold')">
                <a-input-number :value="n(style.fontWeight, 400)" @change="onStyleChange('fontWeight', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item v-if="style.textAlign !== undefined" :label="$t('displayConfig.Properties.ComponentTextAlign')">
                <a-select :value="style.textAlign" @change="onStyleChange('textAlign', $event)">
                  <a-select-option value="left">left</a-select-option>
                  <a-select-option value="center">center</a-select-option>
                  <a-select-option value="right">right</a-select-option>
                  <a-select-option value="justify">justify</a-select-option>
                </a-select>
              </a-form-item>

              <template v-for="(diyData, diyIndex) in styleDiy">
                <a-form-item :key="'diy_' + diyIndex" :label="$t(diyData.name)" v-if="diyData.type != 15">
                  <a-input-number
                    v-if="diyData.type == 1"
                    :value="diyData.value"
                    :min="diyData.min"
                    :max="diyData.max"
                    @change="onDiyChange(diyIndex, $event)"
                    style="width:100%"
                  />
                  <input
                    v-else-if="diyData.type == 2"
                    type="color"
                    :value="normColor(diyData.value)"
                    @input="onDiyChange(diyIndex, $event.target.value)"
                    class="cip"
                  />
                  <a-select v-else-if="diyData.type == 3" :value="diyData.value" @change="onDiyChange(diyIndex, $event)">
                    <a-select-option v-for="option in fontFamilyOptions" :key="option" :value="option">{{ option }}</a-select-option>
                  </a-select>
                  <a-input
                    v-else-if="diyData.type == 4"
                    :value="diyData.value"
                    @change="onDiyInput(diyIndex, $event)"
                  />
                  <div v-else-if="diyData.type == 5" class="image-picker-block">
                    <div class="image-preview" @click="showSystemImageModel(1, diyIndex)">
                      <img v-if="diyData.value" :src="diyData.value" class="image-preview-img" />
                      <div v-else class="image-preview-empty">
                        {{ $t('component.systemImageModel.selectImage') }}
                      </div>
                    </div>
                    <div class="image-picker-actions">
                      <a-button size="small" type="primary" @click="showSystemImageModel(1, diyIndex)">
                        {{ $t('component.systemImageModel.selectImage') }}
                      </a-button>
                      <a-button size="small" @click="onDiyChange(diyIndex, '')">
                        {{ $t('component.systemImageModel.delImage') }}
                      </a-button>
                    </div>
                  </div>
                  <a-select v-else-if="diyData.type == 6" :value="diyData.value" @change="onDiyChange(diyIndex, $event)">
                    <a-select-option v-for="option in diyData.enumList || []" :key="option.value" :value="option.value">
                      {{ $t(option.option) }}
                    </a-select-option>
                  </a-select>
                  <a-input-number
                    v-else-if="diyData.type == 7"
                    :value="diyData.value"
                    :step="0.1"
                    :min="diyData.min"
                    :max="diyData.max"
                    @change="onDiyChange(diyIndex, $event)"
                    style="width:100%"
                  />
                  <a-textarea
                    v-else-if="diyData.type == 9"
                    :value="diyData.value"
                    :rows="diyData.rows || 4"
                    @change="onDiyInput(diyIndex, $event)"
                  />
                  <a-input
                    v-else
                    :value="stringifyDiyValue(diyData.value)"
                    disabled
                  />
                </a-form-item>
              </template>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane
          key="active"
          v-if="Array.isArray(source2D.active)"
          :tab="$t('displayConfig.Properties.TabHeaterActive')"
          :forceRender="true"
        >
          <div class="tab-section" v-for="(activeItem, index) in source2D.active" :key="'active_' + index">
            <div class="tab-section-title">{{ $t(activeItem.name) }}</div>
            <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
              <a-form-item
                v-if="activeItem.condition && activeItem.condition.dataID"
                :label="$t('displayConfig.Properties.ComponentFromBandDevice')"
              >
                <a-checkbox :checked="!!activeItem.condition.isBandDevice" disabled>
                  <span v-if="activeItem.condition.isBandDevice">{{ $t('displayConfig.Properties.ComponentIsBandDevice') }}</span>
                  <span v-else>{{ $t('displayConfig.Properties.ComponentIsBandDeviceModel') }}</span>
                </a-checkbox>
              </a-form-item>
              <a-form-item
                v-if="activeItem.condition && activeItem.condition.isBandDevice"
                :label="$t('displayConfig.Properties.ComponentBandDevice')"
              >
                <a-input :value="activeItem.condition.DeviceName || activeItem.condition.deviceSN || ''" disabled />
              </a-form-item>
              <a-form-item v-if="activeItem.condition" :label="$t('displayConfig.Properties.ComponentBandData')">
                <a-input :value="activeItem.condition.dataName || ''" disabled>
                  <a-tooltip slot="addonAfter" placement="top">
                    <template slot="title">
                      <span>{{ $t('displayConfig.Properties.SelectValue') }}</span>
                    </template>
                    <icon-font @click="showDeviceDataModel(index, 'active')" type="icon-xuanzeshuju" />
                  </a-tooltip>
                </a-input>
              </a-form-item>
              <template v-if="activeItem.isExpression && activeItem.condition">
                <a-form-item :label="$t('component.public.Operator')">
                  <a-select :value="activeItem.condition.operator" @change="onActiveConditionChange(index, 'operator', $event)">
                    <a-select-option value="==">=</a-select-option>
                    <a-select-option value=">">&gt;</a-select-option>
                    <a-select-option value=">=">&gt;=</a-select-option>
                    <a-select-option value="<">&lt;</a-select-option>
                    <a-select-option value="<=">&lt;=</a-select-option>
                    <a-select-option value="<>">&amp;&amp;</a-select-option>
                    <a-select-option value="!=">!=</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('component.public.OperatorValue')">
                  <a-input :value="activeItem.condition.OperatorValue" @change="onActiveConditionInput(index, 'OperatorValue', $event)" />
                </a-form-item>
                <a-form-item
                  v-if="activeItem.condition.operator === '<>'"
                  :label="$t('component.public.OperatorMaxValue')"
                >
                  <a-input :value="activeItem.condition.OperatorMaxValue" @change="onActiveConditionInput(index, 'OperatorMaxValue', $event)" />
                </a-form-item>
              </template>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane
          key="animate"
          v-if="source2D.animate"
          :tab="$t('displayConfig.Properties.TabHeaterAnimate')"
          :forceRender="true"
        >
          <div class="tab-section">
            <div class="tab-section-title">{{ $t('component.public.animate') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
              <a-form-item :label="$t('component.public.animate')">
                <a-select
                  mode="multiple"
                  :value="source2D.animate.selected || []"
                  @change="onAnimateSelectedChange"
                >
                  <a-select-option v-for="item in source2D.animate.animateList || []" :key="item.id" :value="item.id">
                    {{ $t(item.name) }}
                  </a-select-option>
                  <a-select-option value="visible">{{ $t('component.public.animateVisible') }}</a-select-option>
                  <a-select-option value="animateMove">{{ $t('component.public.animateMove') }}</a-select-option>
                </a-select>
              </a-form-item>

              <template v-if="animateSelected.includes('animateMove')">
                <a-form-item :label="$t('component.public.animateX')">
                  <a-input :value="animateMoveX.dataName || ''" disabled>
                    <a-tooltip slot="addonAfter" placement="top">
                      <template slot="title">
                        <span>{{ $t('displayConfig.Properties.SelectValue') }}</span>
                      </template>
                      <icon-font @click="showDeviceDataModel(0, 'animateMoveX')" type="icon-xuanzeshuju" />
                    </a-tooltip>
                  </a-input>
                </a-form-item>
                <a-form-item :label="$t('component.public.animateY')">
                  <a-input :value="animateMoveY.dataName || ''" disabled>
                    <a-tooltip slot="addonAfter" placement="top">
                      <template slot="title">
                        <span>{{ $t('displayConfig.Properties.SelectValue') }}</span>
                      </template>
                      <icon-font @click="showDeviceDataModel(0, 'animateMoveY')" type="icon-xuanzeshuju" />
                    </a-tooltip>
                  </a-input>
                </a-form-item>
              </template>

              <a-form-item
                v-if="source2D.animate.condition && source2D.animate.selected && source2D.animate.selected.length > 0 && !animateSelected.includes('animateMove')"
                :label="$t('component.public.isExpression')"
              >
                <a-radio-group :value="source2D.animate.isExpression" @change="onAnimateExpressionChange">
                  <a-radio :value="true">{{ $t('component.public.Enable') }}</a-radio>
                  <a-radio :value="false">{{ $t('component.public.Forbidden') }}</a-radio>
                </a-radio-group>
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section" v-if="source2D.animate.isExpression && !animateSelected.includes('animateMove')">
            <div class="tab-section-title">{{ $t('displayConfig.Properties.Condition') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
              <a-form-item
                v-if="source2D.animate.condition && source2D.animate.condition.dataID"
                :label="$t('displayConfig.Properties.ComponentFromBandDevice')"
              >
                <a-checkbox :checked="!!source2D.animate.condition.isBandDevice" disabled>
                  <span v-if="source2D.animate.condition.isBandDevice">{{ $t('displayConfig.Properties.ComponentIsBandDevice') }}</span>
                  <span v-else>{{ $t('displayConfig.Properties.ComponentIsBandDeviceModel') }}</span>
                </a-checkbox>
              </a-form-item>
              <a-form-item
                v-if="source2D.animate.condition && source2D.animate.condition.isBandDevice"
                :label="$t('displayConfig.Properties.ComponentBandDevice')"
              >
                <a-input :value="source2D.animate.condition.DeviceName || source2D.animate.condition.deviceSN || ''" disabled />
              </a-form-item>
              <a-form-item v-if="source2D.animate.condition" :label="$t('displayConfig.Properties.ComponentBandData')">
                <a-input :value="source2D.animate.condition.dataName || ''" disabled>
                  <a-tooltip slot="addonAfter" placement="top">
                    <template slot="title">
                      <span>{{ $t('displayConfig.Properties.SelectValue') }}</span>
                    </template>
                    <icon-font @click="showDeviceDataModel(0, 'animate')" type="icon-xuanzeshuju" />
                  </a-tooltip>
                </a-input>
              </a-form-item>
              <a-form-item :label="$t('component.public.Operator')">
                <a-select :value="source2D.animate.condition.operator" @change="onAnimateConditionChange('operator', $event)">
                  <a-select-option value="==">=</a-select-option>
                  <a-select-option value=">">&gt;</a-select-option>
                  <a-select-option value=">=">&gt;=</a-select-option>
                  <a-select-option value="<">&lt;</a-select-option>
                  <a-select-option value="<=">&lt;=</a-select-option>
                  <a-select-option value="<>">&amp;&amp;</a-select-option>
                  <a-select-option value="<!>">||</a-select-option>
                  <a-select-option value="!=">!=</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('component.public.OperatorValue')">
                <a-input :value="source2D.animate.condition.OperatorValue" @change="onAnimateConditionInput('OperatorValue', $event)" />
              </a-form-item>
              <a-form-item
                v-if="source2D.animate.condition.operator === '<>' || source2D.animate.condition.operator === '<!>'"
                :label="$t('component.public.OperatorMaxValue')"
              >
                <a-input :value="source2D.animate.condition.OperatorMaxValue" @change="onAnimateConditionInput('OperatorMaxValue', $event)" />
              </a-form-item>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane
          key="behavior"
          :tab="$t('displayConfig.Properties.TabHeaterBehavior')"
          :forceRender="true"
        >
          <div class="behavior-actions">
            <a-button size="small" type="primary" @click="addBindAction">
              {{ $t('displayConfig.Properties.ComponentAddBandActionBtn') }}
            </a-button>
          </div>
          <div class="tab-section" v-for="(eventItem, index) in source2D.action" :key="'action_' + index">
            <div class="tab-section-title">
              {{ $t('displayConfig.Properties.TabHeaterBehavior') }}-{{ index + 1 }}
              <a-icon type="delete" class="section-delete" @click="delBindAction(index)" />
            </div>
            <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
              <a-form-item :label="$t('displayConfig.Properties.ComponentEvent')">
                <a-select :value="eventItem.type" @change="onActionChange(index, 'type', $event)">
                  <a-select-option v-for="option in actionEventOptions" :key="option.value" :value="option.value">
                    {{ $t(option.label) }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentActionAuth')">
                <a-select mode="multiple" :value="eventItem.actionAuth || []" @change="onActionChange(index, 'actionAuth', $event)">
                  <a-select-option v-for="option in roleOptions" :key="option.value" :value="option.value">
                    {{ $t(option.label) }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentActionPassword')">
                <a-input
                  :value="eventItem.ActionPassword || ''"
                  @change="onActionInput(index, 'ActionPassword', $event)"
                  :style="{ 'text-security': 'disc', '-webkit-text-security': 'disc' }"
                />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentActionVoice')">
                <a-input :value="eventItem.actionVoice" @change="onActionInput(index, 'actionVoice', $event)" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.SecondConfirm')">
                <a-switch :checked="!!eventItem.actionConfirm" @change="onActionChange(index, 'actionConfirm', $event)" />
              </a-form-item>
              <a-form-item :label="$t('displayConfig.Properties.ComponentAction')">
                <a-select :value="eventItem.action" @change="onActionTypeSelected(index, $event)">
                  <a-select-option v-for="option in actionTypeOptions" :key="option.value" :value="option.value">
                    {{ $t(option.label) }}
                  </a-select-option>
                </a-select>
              </a-form-item>

              <template v-if="eventItem.action === 'visible'">
                <a-form-item :label="$t('displayConfig.Properties.ComponentClickShow')">
                  <a-select mode="multiple" :value="eventItem.showItems || []" @change="onActionChange(index, 'showItems', $event)">
                    <a-select-option v-for="option in groupOptions" :key="option.value" :value="option.value">
                      {{ option.label }}
                    </a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.ComponentClickHide')">
                  <a-select mode="multiple" :value="eventItem.hideItems || []" @change="onActionChange(index, 'hideItems', $event)">
                    <a-select-option v-for="option in groupOptions" :key="option.value" :value="option.value">
                      {{ option.label }}
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </template>

              <template v-if="eventItem.action === 'SetValue'">
                <a-form-item :label="$t('displayConfig.Properties.SetDelay')">
                  <a-input :value="eventItem.SetDelay || ''" @change="onActionInput(index, 'SetDelay', $event)" />
                </a-form-item>
                <div class="setvalue-actions">
                  <a-button size="small" type="primary" @click="addBindSetValue(index)">
                    {{ $t('displayConfig.Properties.AddSetValue') }}
                  </a-button>
                </div>
                <div
                  class="setvalue-card"
                  v-for="(setValueItem, setValueIndex) in eventItem.setValue || []"
                  :key="'set_' + index + '_' + setValueIndex"
                >
                  <div class="setvalue-card-title">
                    {{ $t('displayConfig.Properties.SetValueInfo') }}-{{ setValueIndex + 1 }}
                  </div>
                  <a-form-item :label="$t('displayConfig.Properties.IsManualLabel')">
                    <a-checkbox :checked="!!setValueItem.IsManual" @change="onSetValueChecked(index, setValueIndex, 'IsManual', $event)">
                      {{ $t('displayConfig.Properties.IsManual') }}
                    </a-checkbox>
                  </a-form-item>
                  <a-form-item v-if="setValueItem.dataID" :label="$t('displayConfig.Properties.ComponentFromBandDevice')">
                    <a-checkbox :checked="!!setValueItem.isBandDevice" disabled>
                      <span v-if="setValueItem.isBandDevice">{{ $t('displayConfig.Properties.ComponentIsBandDevice') }}</span>
                      <span v-else>{{ $t('displayConfig.Properties.ComponentIsBandDeviceModel') }}</span>
                    </a-checkbox>
                  </a-form-item>
                  <a-form-item v-if="setValueItem.isBandDevice" :label="$t('displayConfig.Properties.ComponentBandDevice')">
                    <a-input :value="setValueItem.DeviceName || setValueItem.deviceSN || ''" disabled />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentBandData')">
                    <a-input :value="setValueItem.dataName || ''" disabled>
                      <a-tooltip slot="addonAfter" placement="top">
                        <template slot="title">
                          <span>{{ $t('displayConfig.Properties.SelectValue') }}</span>
                        </template>
                        <icon-font @click="showSetValueDataModel(index, setValueIndex)" type="icon-xuanzeshuju" />
                      </a-tooltip>
                    </a-input>
                  </a-form-item>
                  <a-form-item v-if="!setValueItem.IsManual" :label="$t('displayConfig.Properties.AutoSetValue')">
                    <a-input :value="setValueItem.AutoSetValue" @change="onSetValueInput(index, setValueIndex, 'AutoSetValue', $event)" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.SetPassword')">
                    <a-input :value="setValueItem.SetPassword" @change="onSetValueInput(index, setValueIndex, 'SetPassword', $event)" />
                  </a-form-item>
                </div>
              </template>

              <template v-if="eventItem.action === 'link'">
                <a-form-item :label="$t('displayConfig.Properties.linkType')">
                  <a-select :value="eventItem.link && eventItem.link.linkType" @change="onActionLinkChange(index, 'linkType', $event)">
                    <a-select-option value="Inside">{{ $t('displayConfig.Properties.linkInside') }}</a-select-option>
                    <a-select-option value="External">{{ $t('displayConfig.Properties.linkExternal') }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item v-if="!eventItem.link || eventItem.link.linkType !== 'Inside'" :label="$t('displayConfig.Properties.linkExternalUrl')">
                  <a-input :value="eventItem.link && eventItem.link.External" @change="onActionLinkInput(index, 'External', $event)" />
                </a-form-item>
                <a-form-item v-if="!eventItem.link || eventItem.link.linkType !== 'Inside'" :label="$t('displayConfig.Properties.OpenLinkExternalType')">
                  <a-select :value="(eventItem.link && eventItem.link.OpenExternalType) || 'new'" @change="onActionLinkChange(index, 'OpenExternalType', $event)">
                    <a-select-option value="self">{{ $t('displayConfig.Properties.OpenLinkExternalSelf') }}</a-select-option>
                    <a-select-option value="new">{{ $t('displayConfig.Properties.OpenLinkExternalNew') }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.isLinkPopUp')">
                  <a-checkbox :checked="!!(eventItem.link && eventItem.link.isPopUp)" @change="onActionLinkChecked(index, 'isPopUp', $event)">
                    {{ $t('displayConfig.Properties.isLinkPopUp') }}
                  </a-checkbox>
                </a-form-item>
                <a-form-item v-if="eventItem.link && eventItem.link.isPopUp" :label="$t('displayConfig.Properties.autoClose')">
                  <a-checkbox :checked="!!eventItem.link.autoClose" @change="onActionLinkChecked(index, 'autoClose', $event)">
                    {{ $t('displayConfig.Properties.autoClose') }}
                  </a-checkbox>
                </a-form-item>
                <template v-if="eventItem.link && eventItem.link.isPopUp && eventItem.link.linkType !== 'Inside'">
                  <a-form-item :label="$t('displayConfig.Properties.linkExternalWidth')">
                    <a-input :value="eventItem.link.width || ''" @change="onActionLinkInput(index, 'width', $event)" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.linkExternalHeight')">
                    <a-input :value="eventItem.link.height || ''" @change="onActionLinkInput(index, 'height', $event)" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.linkExternalTitle')">
                    <a-input :value="eventItem.link.title || ''" @change="onActionLinkInput(index, 'title', $event)" />
                  </a-form-item>
                </template>
              </template>

              <template v-if="eventItem.action === 'SysScript'">
                <a-form-item :label="$t('displayConfig.Properties.action.SysScript')">
                  <a-select mode="multiple" :value="eventItem.ScriptList || []" @change="onActionChange(index, 'ScriptList', $event)">
                    <a-select-option v-for="script in scriptOptions" :key="script.value" :value="script.value">
                      {{ script.label }}
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </template>

              <template v-if="eventItem.action === 'Animation'">
                <a-form-item :label="$t('displayConfig.Properties.action.animationList')">
                  <a-select mode="multiple" :value="animateSelected" @change="onAnimateSelectedChange">
                    <a-select-option v-for="item in source2D.animate && source2D.animate.animateList ? source2D.animate.animateList : []" :key="item.id" :value="item.id">
                      {{ $t(item.name) }}
                    </a-select-option>
                    <a-select-option value="visible">{{ $t('component.public.animateVisible') }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.action.animationStatus')">
                  <a-select :value="eventItem.animationStatus" @change="onActionChange(index, 'animationStatus', $event)">
                    <a-select-option value="start">{{ $t('displayConfig.Properties.action.animationStart') }}</a-select-option>
                    <a-select-option value="stop">{{ $t('displayConfig.Properties.action.animationStop') }}</a-select-option>
                  </a-select>
                </a-form-item>
              </template>

              <template v-if="eventItem.action === 'RestApi'">
                <a-form-item :label="$t('displayConfig.Properties.RestApiName')">
                  <a-input :value="eventItem.RestApi && eventItem.RestApi.Name" @change="onActionNestedInput(index, 'RestApi', 'Name', $event)" />
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.RestApiFrom')">
                  <a-select :value="(eventItem.RestApi && eventItem.RestApi.IsSystem) || '1'" @change="onActionNestedChange(index, 'RestApi', 'IsSystem', $event)">
                    <a-select-option value="1">{{ $t('displayConfig.Properties.ExternRestApi') }}</a-select-option>
                    <a-select-option value="2">{{ $t('displayConfig.Properties.SystemRestApi') }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.RestApiType')">
                  <a-select :value="(eventItem.RestApi && eventItem.RestApi.Type) || 'Post'" @change="onActionNestedChange(index, 'RestApi', 'Type', $event)">
                    <a-select-option value="Post">Post</a-select-option>
                    <a-select-option value="Get">Get</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.RestApiUrl')">
                  <a-input :value="eventItem.RestApi && eventItem.RestApi.Url" @change="onActionNestedInput(index, 'RestApi', 'Url', $event)" />
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.RestApiParam')">
                  <a-textarea :rows="6" :value="eventItem.RestApi && eventItem.RestApi.Params" @change="onActionNestedInput(index, 'RestApi', 'Params', $event)" />
                </a-form-item>
              </template>
            </a-form>
          </div>
        </a-tab-pane>
      </a-tabs>
    </template>

    <device-data-model ref="deviceDataModel" @onSelectDataModel="onSelectData"></device-data-model>
    <system-image-model
      ref="systemImageModel"
      :networkImageUrl="imageComponentData"
      @onSelectImage="onSelectImage"
    ></system-image-model>
  </div>
</template>

<script>
import deviceDataModel from '@/components/deviceDataModel/deviceDataModel.vue'
import systemImageModel from '@/components/systemImageModel/systemImageModel.vue'
import { GetScriptList } from '@/services/ismscripts'

export default {
  name: 'Properties2DPanel',
  i18n: require('@/i18n/language'),
  components: {
    deviceDataModel,
    systemImageModel
  },
  props: {
    currentObject: { type: Object, default: null }
  },
  data() {
    return {
      activeKey: 'style',
      bandType: '',
      selectBandDataIndex: 0,
      actionIndex: 0,
      imageComponentData: '',
      selectImageType: 1,
      selectBandDiyDataIndex: 0,
      scriptOptions: []
    }
  },
  computed: {
    source2D() {
      return this.currentObject && this.currentObject.source2D ? this.currentObject.source2D : {}
    },
    style() {
      return this.source2D.style || {}
    },
    pos() {
      return this.style.position || {}
    },
    styleDiy() {
      return Array.isArray(this.style.diy) ? this.style.diy : []
    },
    animateSelected() {
      return this.source2D.animate && Array.isArray(this.source2D.animate.selected) ? this.source2D.animate.selected : []
    },
    animateMoveX() {
      return this.source2D.animate && this.source2D.animate.move && this.source2D.animate.move.x ? this.source2D.animate.move.x : {}
    },
    animateMoveY() {
      return this.source2D.animate && this.source2D.animate.move && this.source2D.animate.move.y ? this.source2D.animate.move.y : {}
    },
    actionEventOptions() {
      return [
        { label: 'displayConfig.Properties.action.click', value: 'click' },
        { label: 'displayConfig.Properties.action.MouseDown', value: 'mousedown' },
        { label: 'displayConfig.Properties.action.MouseUp', value: 'mouseup' },
        { label: 'displayConfig.Properties.action.dbClick', value: 'dblclick' },
        { label: 'displayConfig.Properties.action.mouseenter', value: 'mouseenter' },
        { label: 'displayConfig.Properties.action.mouseleave', value: 'mouseleave' }
      ]
    },
    actionTypeOptions() {
      return [
        { label: 'displayConfig.Properties.action.openLink', value: 'link' },
        { label: 'displayConfig.Properties.action.SetValue', value: 'SetValue' },
        { label: 'displayConfig.Properties.action.Visible', value: 'visible' },
        { label: 'displayConfig.Properties.action.SysScript', value: 'SysScript' },
        { label: 'displayConfig.Properties.action.DeviceView', value: 'DeviceView' },
        { label: 'displayConfig.Properties.action.RestApi', value: 'RestApi' },
        { label: 'displayConfig.Properties.action.animation', value: 'Animation' }
      ]
    },
    roleOptions() {
      return [
        { label: 'account.settings.UserList.RoleOperator', value: 'Operator' },
        { label: 'account.settings.UserList.RoleUser', value: 'User' }
      ]
    },
    groupOptions() {
      const activeList = Array.isArray(this.source2D.active) ? this.source2D.active : []
      return activeList.map((item, index) => ({
        label: this.$t(item.name || ('Group ' + (index + 1))),
        value: item.id || ('active_' + index)
      }))
    }
  },
  methods: {
    stringifyDiyValue(val) {
      if (val === null || val === undefined) return ''
      if (typeof val === 'object') return JSON.stringify(val)
      return String(val)
    },
    n(val, fallback) {
      const num = parseFloat(val)
      return isNaN(num) ? fallback : num
    },
    normColor(val) {
      if (!val || typeof val !== 'string' || val === 'transparent') return '#000000'
      if (/^#[0-9a-fA-F]{3}$/.test(val)) {
        return '#' + val[1] + val[1] + val[2] + val[2] + val[3] + val[3]
      }
      if (/^#[0-9a-fA-F]{6}$/.test(val)) return val.toLowerCase()
      if (/^#[0-9a-fA-F]{8}$/.test(val)) return val.slice(0, 7).toLowerCase()
      return '#000000'
    },
    ensureStyle() {
      if (!this.currentObject.source2D) this.$set(this.currentObject, 'source2D', {})
      if (!this.currentObject.source2D.style) this.$set(this.currentObject.source2D, 'style', {})
      if (!this.currentObject.source2D.style.position) this.$set(this.currentObject.source2D.style, 'position', {})
      if (!Array.isArray(this.currentObject.source2D.style.diy)) this.$set(this.currentObject.source2D.style, 'diy', [])
    },
    ensureActive(index) {
      if (!Array.isArray(this.currentObject.source2D.active)) {
        this.$set(this.currentObject.source2D, 'active', [])
      }
      if (!this.currentObject.source2D.active[index].condition) {
        this.$set(this.currentObject.source2D.active[index], 'condition', {})
      }
    },
    ensureAnimate() {
      if (!this.currentObject.source2D.animate) this.$set(this.currentObject.source2D, 'animate', { selected: [], condition: {}, move: { x: {}, y: {} } })
      if (!Array.isArray(this.currentObject.source2D.animate.selected)) this.$set(this.currentObject.source2D.animate, 'selected', [])
      if (!this.currentObject.source2D.animate.condition) this.$set(this.currentObject.source2D.animate, 'condition', {})
      if (!this.currentObject.source2D.animate.move) this.$set(this.currentObject.source2D.animate, 'move', { x: {}, y: {} })
      if (!this.currentObject.source2D.animate.move.x) this.$set(this.currentObject.source2D.animate.move, 'x', {})
      if (!this.currentObject.source2D.animate.move.y) this.$set(this.currentObject.source2D.animate.move, 'y', {})
    },
    ensureActions() {
      if (!Array.isArray(this.currentObject.source2D.action)) {
        this.$set(this.currentObject.source2D, 'action', [])
      }
    },
    ensureActionLink(index) {
      this.ensureActions()
      if (!this.currentObject.source2D.action[index].link) {
        this.$set(this.currentObject.source2D.action[index], 'link', {
          linkType: 'Inside',
          Inside: { displayUUID: '', pageUUID: '' },
          isPopUp: false,
          autoClose: false,
          External: '',
          OpenExternalType: 'new',
          width: '',
          height: '',
          title: ''
        })
      }
    },
    ensureSetValue(index) {
      this.ensureActions()
      if (!Array.isArray(this.currentObject.source2D.action[index].setValue)) {
        this.$set(this.currentObject.source2D.action[index], 'setValue', [])
      }
    },
    ensureActionDefaults(index) {
      this.ensureActions()
      const target = this.currentObject.source2D.action[index]
      if (!target) return
      if (!Array.isArray(target.actionAuth)) this.$set(target, 'actionAuth', [])
      if (target.ActionPassword === undefined) this.$set(target, 'ActionPassword', '')
      if (target.actionVoice === undefined) this.$set(target, 'actionVoice', '')
      if (target.actionConfirm === undefined) this.$set(target, 'actionConfirm', false)
      if (!Array.isArray(target.showItems)) this.$set(target, 'showItems', [])
      if (!Array.isArray(target.hideItems)) this.$set(target, 'hideItems', [])
      this.ensureActionLink(index)
      if (!target.RestApi) {
        this.$set(target, 'RestApi', { Name: '', IsSystem: '1', Type: 'Post', Url: '', Params: '{}' })
      }
      if (!target.DeviceView) {
        this.$set(target, 'DeviceView', { key: '', showUUID: '', showPageUUID: '', type: '', selectKey: '', isPopUp: false })
      }
      if (!Array.isArray(target.ScriptList)) this.$set(target, 'ScriptList', [])
      if (target.animationStatus === undefined) this.$set(target, 'animationStatus', '')
      if (target.action === 'SetValue') {
        this.ensureSetValue(index)
        if (!target.setValue.length) {
          target.setValue.push({
            deviceSN: '',
            DeviceName: '',
            IsManual: false,
            AutoSetValue: '',
            SetPassword: '',
            isBandDevice: false,
            dataID: '',
            dataName: ''
          })
        }
      }
    },
    onRootPropChange(prop, e) {
      this.$set(this.currentObject, prop, e.target.value)
      this.$emit('prop-change')
    },
    onStyleChange(prop, value) {
      this.ensureStyle()
      this.$set(this.currentObject.source2D.style, prop, value)
      this.$emit('prop-change')
    },
    onStyleInput(prop, e) {
      this.onStyleChange(prop, e.target.value)
    },
    onStyleSwitch(prop, value) {
      this.onStyleChange(prop, value)
    },
    onPositionChange(prop, value) {
      this.ensureStyle()
      this.$set(this.currentObject.source2D.style.position, prop, value)
      this.$emit('prop-change')
    },
    onDiyChange(index, value) {
      this.ensureStyle()
      if (!this.currentObject.source2D.style.diy[index]) return
      this.$set(this.currentObject.source2D.style.diy[index], 'value', value)
      this.$emit('prop-change')
    },
    onDiyInput(index, e) {
      this.onDiyChange(index, e.target.value)
    },
    onActiveConditionChange(index, prop, value) {
      this.ensureActive(index)
      this.$set(this.currentObject.source2D.active[index].condition, prop, value)
      this.$emit('prop-change')
    },
    onActiveConditionInput(index, prop, e) {
      this.onActiveConditionChange(index, prop, e.target.value)
    },
    onAnimateSelectedChange(value) {
      this.ensureAnimate()
      this.$set(this.currentObject.source2D.animate, 'selected', value)
      this.$emit('prop-change')
    },
    onAnimateExpressionChange(e) {
      this.ensureAnimate()
      this.$set(this.currentObject.source2D.animate, 'isExpression', e.target.value)
      this.$emit('prop-change')
    },
    onAnimateConditionChange(prop, value) {
      this.ensureAnimate()
      this.$set(this.currentObject.source2D.animate.condition, prop, value)
      this.$emit('prop-change')
    },
    onAnimateConditionInput(prop, e) {
      this.onAnimateConditionChange(prop, e.target.value)
    },
    onActionChange(index, prop, value) {
      this.ensureActions()
      this.$set(this.currentObject.source2D.action[index], prop, value)
      this.ensureActionDefaults(index)
      this.$emit('prop-change')
    },
    onActionTypeSelected(index, value) {
      this.onActionChange(index, 'action', value)
      if (value === 'Animation') {
        this.ensureAnimate()
        this.$set(this.currentObject.source2D.animate, 'isExpression', true)
      }
    },
    onActionInput(index, prop, e) {
      this.onActionChange(index, prop, e.target.value)
    },
    onActionNestedChange(index, parentProp, prop, value) {
      this.ensureActions()
      this.ensureActionDefaults(index)
      const parent = this.currentObject.source2D.action[index][parentProp]
      if (!parent) return
      this.$set(parent, prop, value)
      this.$emit('prop-change')
    },
    onActionNestedInput(index, parentProp, prop, e) {
      this.onActionNestedChange(index, parentProp, prop, e.target.value)
    },
    onActionLinkChange(index, prop, value) {
      this.ensureActionLink(index)
      this.$set(this.currentObject.source2D.action[index].link, prop, value)
      this.$emit('prop-change')
    },
    onActionLinkInput(index, prop, e) {
      this.onActionLinkChange(index, prop, e.target.value)
    },
    onActionLinkChecked(index, prop, e) {
      this.onActionLinkChange(index, prop, e.target.checked)
    },
    addBindSetValue(index) {
      this.ensureSetValue(index)
      this.currentObject.source2D.action[index].setValue.push({
        deviceSN: '',
        DeviceName: '',
        IsManual: false,
        AutoSetValue: '',
        SetPassword: '',
        isBandDevice: false,
        dataID: '',
        dataName: ''
      })
      this.$emit('prop-change')
    },
    addBindAction() {
      this.ensureActions()
      this.currentObject.source2D.action.push({
        type: 'click',
        action: 'SetValue',
        actionAuth: [],
        ActionPassword: '',
        actionVoice: '',
        actionConfirm: false,
        showItems: [],
        hideItems: [],
        link: {
          linkType: 'Inside',
          Inside: { displayUUID: '', pageUUID: '' },
          isPopUp: false,
          autoClose: false,
          External: '',
          OpenExternalType: 'new',
          width: '',
          height: '',
          title: ''
        },
        setValue: [
          {
            deviceSN: '',
            DeviceName: '',
            IsManual: false,
            AutoSetValue: '',
            SetPassword: '',
            isBandDevice: false,
            dataID: '',
            dataName: ''
          }
        ],
        RestApi: { Name: '', IsSystem: '1', Type: 'Post', Url: '', Params: '{}' },
        DeviceView: { key: '', showUUID: '', showPageUUID: '', type: '', selectKey: '', isPopUp: false },
        ScriptList: [],
        animationStatus: ''
      })
      this.$emit('prop-change')
    },
    delBindAction(index) {
      this.ensureActions()
      this.currentObject.source2D.action.splice(index, 1)
      this.$emit('prop-change')
    },
    onSetValueInput(actionIndex, setValueIndex, prop, e) {
      this.ensureSetValue(actionIndex)
      this.$set(this.currentObject.source2D.action[actionIndex].setValue[setValueIndex], prop, e.target.value)
      this.$emit('prop-change')
    },
    onSetValueChecked(actionIndex, setValueIndex, prop, e) {
      this.ensureSetValue(actionIndex)
      this.$set(this.currentObject.source2D.action[actionIndex].setValue[setValueIndex], prop, e.target.checked)
      this.$emit('prop-change')
    },
    showSetValueDataModel(actionIndex, setValueIndex) {
      this.selectBandDataIndex = setValueIndex
      this.bandType = 'setValue'
      this.actionIndex = actionIndex
      if (this.$refs.deviceDataModel) {
        this.$refs.deviceDataModel.showDataModal()
      }
    },
    showDeviceDataModel(index, type) {
      this.selectBandDataIndex = index
      this.bandType = type
      if (this.$refs.deviceDataModel) {
        this.$refs.deviceDataModel.showDataModal()
      }
    },
    showSystemImageModel(selectType, bindIndex, showType) {
      this.selectImageType = selectType
      this.selectBandDiyDataIndex = bindIndex
      if (this.selectImageType === 1) {
        this.imageComponentData = (this.styleDiy[bindIndex] && this.styleDiy[bindIndex].value) || ''
      } else {
        this.imageComponentData = ''
      }
      if (this.$refs.systemImageModel) {
        this.$refs.systemImageModel.showModal(showType)
      }
    },
    onSelectImage(url) {
      if (this.selectImageType === 1) {
        this.onDiyChange(this.selectBandDiyDataIndex, url)
      }
    },
    onSelectData(selectData) {
      if (this.bandType === 'active') {
        this.ensureActive(this.selectBandDataIndex)
        const target = this.currentObject.source2D.active[this.selectBandDataIndex].condition
        this.$set(target, 'isBandDevice', selectData.IsDevice)
        this.$set(target, 'deviceSN', selectData.DeviceSN)
        this.$set(target, 'dataName', selectData.name)
        this.$set(target, 'dataID', selectData.uuid)
        this.$set(target, 'dataUnit', selectData.unit)
        this.$set(target, 'DeviceName', selectData.DeviceName)
      } else if (this.bandType === 'animate') {
        this.ensureAnimate()
        const target = this.currentObject.source2D.animate.condition
        this.$set(target, 'isBandDevice', selectData.IsDevice)
        this.$set(target, 'deviceSN', selectData.DeviceSN)
        this.$set(target, 'dataName', selectData.name)
        this.$set(target, 'dataID', selectData.uuid)
        this.$set(target, 'DeviceName', selectData.DeviceName)
      } else if (this.bandType === 'animateMoveX' || this.bandType === 'animateMoveY') {
        this.ensureAnimate()
        const axis = this.bandType === 'animateMoveX' ? 'x' : 'y'
        const target = this.currentObject.source2D.animate.move[axis]
        this.$set(target, 'isBandDevice', selectData.IsDevice)
        this.$set(target, 'deviceSN', selectData.DeviceSN)
        this.$set(target, 'dataName', selectData.name)
        this.$set(target, 'dataID', selectData.uuid)
        this.$set(target, 'DeviceName', selectData.DeviceName)
      } else if (this.bandType === 'setValue') {
        this.ensureSetValue(this.actionIndex)
        const target = this.currentObject.source2D.action[this.actionIndex].setValue[this.selectBandDataIndex]
        this.$set(target, 'isBandDevice', selectData.IsDevice)
        this.$set(target, 'deviceSN', selectData.DeviceSN)
        this.$set(target, 'dataName', selectData.name)
        this.$set(target, 'dataID', selectData.uuid)
        this.$set(target, 'DeviceName', selectData.DeviceName)
      }
      this.$emit('prop-change')
    },
    loadScriptOptions() {
      GetScriptList().then((res) => {
        const list = res && res.data && Array.isArray(res.data.list) ? res.data.list : []
        this.scriptOptions = list.map(item => ({
          value: item.ScriptUuid,
          label: item.ScriptName
        }))
      }).catch(() => {
        this.scriptOptions = []
      })
    }
  },
  mounted() {
    this.loadScriptOptions()
  }
}
</script>

<style scoped>
.right-panel { width:100%;height:100%;background:#fafafa;border-left:1px solid #e8e8e8;display:flex;flex-direction:column;overflow-y:auto; }
.no-select-tip { padding:24px 16px;text-align:center;color:#bbb;font-size:13px; }
.no-select-tip i { font-size:32px;margin-bottom:10px;display:block; }
.prop-tabs { flex:1;display:flex;flex-direction:column;min-height:0;overflow-y:auto; }
.tab-section { margin-bottom:8px;background:#fff;border:1px solid #e8e8e8;border-radius:4px; }
.tab-section-title { padding:8px 12px;font-size:12px;font-weight:600;color:#666;background:#f5f5f5;border-bottom:1px solid #e8e8e8;letter-spacing:0.6px; }
.cip { width:36px;height:26px;border:1px solid #3d3d7a;border-radius:4px;cursor:pointer;background:transparent;padding:1px; }
.behavior-actions { margin: 0 0 8px; text-align: center; }
.setvalue-actions { margin: 8px 0 12px; text-align: center; }
.setvalue-card { margin: 8px 0; padding: 8px; border: 1px solid #d9e6f2; border-radius: 4px; background: #fafcff; }
.setvalue-card-title { margin-bottom: 8px; font-size: 12px; font-weight: 600; color: #666; }
.section-delete { float:right; cursor:pointer; color:#ff4d4f; }
.image-picker-block { display:flex; flex-direction:column; gap:8px; }
.image-preview { width:100%; min-height:120px; border:1px solid #d9d9d9; border-radius:4px; background:#fafafa; display:flex; align-items:center; justify-content:center; cursor:pointer; overflow:hidden; }
.image-preview-img { display:block; width:100%; max-height:220px; object-fit:contain; }
.image-preview-empty { padding:16px; color:#999; font-size:12px; text-align:center; }
.image-picker-actions { display:flex; gap:8px; }
</style>
