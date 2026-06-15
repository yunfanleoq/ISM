# Industrial 3D Assets

Put optimized `.glb` files for the ISM 3D editor catalog in this folder.

The current catalog in `src/pages/ISM3DEditor/assets/industrialAssets.js` expects these paths:

- `platform/pipe-bridge.glb`
- `platform/maintenance-platform.glb`
- `vessel/distillation-tower.glb`
- `vessel/absorber-column.glb`
- `separation/plate-filter.glb`
- `separation/bag-filter.glb`
- `pump/centrifugal-pump.glb`
- `electrical/control-cabinet.glb`
- `electrical/mcc-panel.glb`
- `robot/robot-arm.glb`
- `conveyor/roller-conveyor.glb`
- `warehouse/warehouse-rack.glb`

If a file is missing, the editor creates a low-poly fallback component first, so the asset remains usable.

Recommended model sources to review before adding files:

- Khronos glTF sample assets: https://github.khronos.org/glTF-Assets/
- Sketchfab industrial equipment assets: https://sketchfab.com/tags/industrial-equipment
- Free3D industrial models: https://free3d.com/3d-models/industrial
- TurboSquid free industrial glTF models: https://www.turbosquid.com/3d-model/free/industrial/gltf
- Industrial City GLB Pack CC0: https://eclair-assets.itch.io/industrial-city-glb-pack-25-free-cc0-3d-models

Check each asset's license before committing it to the project.
