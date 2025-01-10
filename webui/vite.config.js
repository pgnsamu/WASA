import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, ssrBuild }) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
	};
	ret.define = {
		// TODO: forse da rimuovere
		// Define Vue feature flags for production builds
		__VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true',
		__VUE_OPTIONS_API__: 'true',
		__VUE_COMPOSITION_API__: 'true',
		// Do not modify this constant, it is used in the evaluation.
		"__API_URL__": JSON.stringify("http://localhost:3000"),
	};
	return ret;
})
