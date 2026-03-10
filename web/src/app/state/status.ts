import { derived, type Readable } from 'svelte/store';
import { state, APP_CONNECTION_STATUS } from '.';
import { status as mavbakeStatus } from './mavbake';
import { status as mavpayStatus } from './mavpay';

export const APP_STATUS_LEVEL = derived(
	[state, APP_CONNECTION_STATUS, mavbakeStatus, mavpayStatus],
	([$state, $connectionStatus, $mavbakeStatus, $mavpayStatus]) => {
		if ($connectionStatus !== 'connected') {
			return 'error';
		}

		for (const node of Object.values($state.nodes)) {
			if (node.connection_status !== 'connected' && node.is_essential) {
				return 'error';
			}
		}
		if ($mavbakeStatus !== 'ok') {
			return $mavbakeStatus;
		}

		if ($mavpayStatus !== 'ok') {
			return $mavpayStatus;
		}

		return 'ok';
	}
) as Readable<'ok' | 'error' | 'warning'>;
