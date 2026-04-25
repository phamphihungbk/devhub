import type { PiniaPluginContext } from 'pinia'

export function resetSetupStore({ store }: PiniaPluginContext) {
  const initialState = structuredClone(store.$state)

  store.$reset = () => {
    store.$patch(structuredClone(initialState))
  }
}
