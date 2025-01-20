import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import { watch } from 'node:fs';

// https://vitejs.dev/config/
export default defineConfig(({command, mode, ssrBuild}) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
		// hmr config
		// TODO se al prof non va bene che ci sia, levarlo.
		server: {
			watch: {
			  usePolling: true, // Use polling if file changes are not detected
			},
			hmr: {
			  overlay: true, // Enable HMR overlay for errors
			},
		  },
	};
	ret.define = {
		// Do not modify this constant, it is used in the evaluation.
		"__API_URL__": JSON.stringify("http://localhost:3000"),
	};
	return ret;
})
