import { EmptyMavpayInfo, type PayoutBlueprint, type MavpayInfo } from "@src/common/types/mavpay"
import { readBody } from "@src/util/fetch"

export async function getMavpayInfo() {
	try {
		const response = await fetch('/api/mavpay/info', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		})

		if (response.status !== 200) {
			throw new Error('Failed to get mavpay info')
		}

		return await response.json() as MavpayInfo
	} catch (e) {
		console.log(e)
		return EmptyMavpayInfo
	}
}

export async function generatePayuts(cycle: number | undefined, cb: (message: string) => void, dry?: boolean) {
	const cycleQuery = cycle ? `cycle=${cycle}` : ''
	const dryQuery = dry ? `dry=${dry}` : ''

	const response = await fetch(`/api/mavpay/generate-payouts?${cycleQuery}&${dryQuery}`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response.status !== 200) {
		throw new Error(response.statusText)
	}

	return await readBody(response, cb)
}

export async function executePayuts(blueprint: PayoutBlueprint, cb: (message: string) => void, dry?: boolean) {
	const dryQuery = dry ? `dry=${dry}` : ''

	const response = await fetch(`/api/mavpay/pay?${dryQuery}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(blueprint)
	})
	if (response.status !== 200) {
		throw new Error(response.statusText)
	}

	const reader = response.body?.getReader();
	if (!reader) {
		throw new Error('No reader')
	}
	const decoder = new TextDecoder('utf-8');
	let buffer = '';


	for (; ;) {
		const { done, value } = await reader.read();
		if (done) {
			if (buffer.length > 0) {
				console.log(`Received chunk: ${buffer}`);
			}
			break;
		}

		buffer += decoder.decode(value, { stream: true });

		let newlineIndex;
		while ((newlineIndex = buffer.indexOf('\n')) > -1) {
			const line = buffer.slice(0, newlineIndex + 1).trim();
			buffer = buffer.slice(newlineIndex + 1);
			cb(line);
		}
	}
}

export async function stopContinual() {
	const response = await fetch('/api/mavpay/stop-continual', {
		method: 'GET',
	})

	if (response.status !== 200) {
		const body = await response.text()
		throw new Error(body)
	}
}

export async function startContinual() {
	const response = await fetch('/api/mavpay/start-continual', {
		method: 'GET',
	})


	if (response.status !== 200) {
		const body = await response.text()
		throw new Error(body)
	}
}

export async function disableContinual() {
	const response = await fetch('/api/mavpay/disable-continual', {
		method: 'GET',
	})

	if (response.status !== 200) {
		const body = await response.text()
		throw new Error(body)
	}
}

export async function enableContinual() {
	const response = await fetch('/api/mavpay/enable-continual', {
		method: 'GET',
	})

	if (response.status !== 200) {
		const body = await response.text()
		throw new Error(body)
	}
}

export async function listReports(dry?: boolean) {
	const response = await fetch(`/api/mavpay/list-reports?dry=${dry === true}`, {
		method: 'GET',
	})

	if (response.status !== 200) {
		throw new Error('Failed to list reports')
	}

	return await response.json() as Array<string>
}

export async function getReport(report: string, dry?: boolean) {
	const response = await fetch(`/api/mavpay/report?id=${report}&dry=${dry === true}`, {
		method: 'GET',
	})

	if (response.status !== 200) {
		throw new Error('Failed to get report')
	}

	return await response.json()
}

export async function testNotify(notificator = 'all', cb: (message: string) => void) {
	const query = notificator === "all" ? "" : `notificator=${notificator}`
	const response = await fetch(`/api/mavpay/test-notify?${query}`, {
		method: 'POST',
	})

	if (response.status !== 200) {
		throw new Error(response.statusText)
	}

	return await readBody(response, cb)
}

export async function testExtensions(cb: (message: string) => void) {
	const response = await fetch(`/api/mavpay/test-extensions`, {
		method: 'POST',
	})

	if (response.status !== 200) {
		throw new Error(response.statusText)
	}

	return await readBody(response, cb)
}