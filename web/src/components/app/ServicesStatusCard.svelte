<script lang="ts">
	import Card from '@components/starlight/components/Card.svelte';
	import type { AmiServiceInfo, ServicesStatus } from '@src/common/types/status';
	import { formatTimestamp, formatTimestampAgo } from '@src/util/format';
	import { onDestroy } from 'svelte';

	export let title = 'Validator Services';
	export let services: ServicesStatus = { timestamp: 0, applications: {} };
	export let hideApplicationTitle = false;

	type ServiceInfo = AmiServiceInfo & { formattedTimestamp: string };

	const interval = setInterval(() => {
		for (const [_, app] of Object.entries(services.applications ?? {})) {
			for (const [_, v] of Object.entries(app)) {
				v.formattedTimestamp = formatTimestampAgo(v.started);
			}
		}
	}, 500);

	$: iterableServices = Object.entries(services.applications ?? {}).map(([id, appServices]) => {
		return [id, Object.entries(appServices)] as [string, Array<[string, ServiceInfo]>];
	});

	onDestroy(() => clearInterval(interval));
</script>

<div class="services-wrap">
	<Card>
		<div class="services">
			<div class="card-title">
				<span class="icon">
					<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
				</span>
				{title}
			</div>

			<div class="services-list">
				<div class="list-header">
					<span>Service</span>
					<span>Status</span>
					<span>Active</span>
				</div>

				{#each iterableServices as [id, appServices]}
					{#if !hideApplicationTitle && iterableServices.length > 1 && appServices.length > 0}
						<div class="app-title">{id}</div>
					{/if}
					{#each appServices as [serviceId, serviceInfo]}
						<div class="service-row">
							<span class="service-name">{serviceId}</span>
							{#if serviceInfo?.status === 'running'}
								<span class="badge badge-active">
									<span class="badge-dot"></span>
									Active
								</span>
							{:else}
								<span class="badge badge-inactive">
									<span class="badge-dot"></span>
									N/A
								</span>
							{/if}
							<span class="service-active">
								{serviceInfo?.status === 'running' ? formatTimestampAgo(serviceInfo.started) : '-'}
							</span>
						</div>
					{:else}
						<div class="service-row">
							<span class="service-app-name">{id}</span>
							<span class="badge badge-inactive">
								<span class="badge-dot"></span>
								N/A
							</span>
							<span class="service-active">-</span>
						</div>
					{/each}
				{/each}
			</div>

			<div class="timestamp">{formatTimestamp(services.timestamp)}</div>
		</div>
	</Card>
</div>

<style lang="sass">
.services-wrap
	height: 100%
	display: grid

	:global(.card)
		height: 100%
		box-sizing: border-box
		--card-accent: linear-gradient(90deg, var(--cyan), var(--cyan-dim))
		--card-hover-shadow: 0 8px 32px var(--cyan-glow)

.services
	display: flex
	flex-direction: column
	height: 100%
	gap: var(--spacing)

.card-title
	display: flex
	align-items: center
	gap: 8px
	font-size: 13px
	font-weight: 600
	text-transform: uppercase
	letter-spacing: 1.5px
	color: var(--cyan)

	.icon
		width: 20px
		height: 20px
		border-radius: 6px
		display: flex
		align-items: center
		justify-content: center
		background: var(--cyan-glow)
		border: 1px solid rgba(34, 211, 238, 0.3)
		flex-shrink: 0

.services-list
	flex: 1
	display: flex
	flex-direction: column
	gap: 4px

	.list-header
		display: grid
		grid-template-columns: 1fr 90px 90px
		gap: 16px
		padding: 0 12px 8px
		border-bottom: 1px solid var(--border-default)

		span
			font-size: 11px
			font-weight: 600
			text-transform: uppercase
			letter-spacing: 1px
			color: var(--text-muted)

		span:not(:first-child)
			text-align: center

	.app-title
		font-size: 11px
		font-weight: 600
		text-transform: uppercase
		letter-spacing: 1px
		color: var(--text-muted)
		padding: 8px 12px 4px

	.service-row
		display: grid
		grid-template-columns: 1fr 90px 90px
		gap: 16px
		align-items: center
		padding: 10px 12px
		border-radius: 4px
		transition: background 0.15s

		&:hover
			background: rgba(255,255,255,0.03)

	.service-name
		font-size: 14px
		font-weight: 500
		color: var(--text-primary)

	.service-app-name
		font-size: 11px
		font-weight: 600
		text-transform: uppercase
		letter-spacing: 1px
		color: var(--text-muted)

	.service-active
		font-family: var(--font-mono)
		font-size: 12px
		color: var(--text-muted)
		text-align: center

.timestamp
	font-family: var(--font-mono)
	font-size: 12px
	color: var(--text-muted)
	padding-top: var(--spacing)
	border-top: 1px solid var(--border-default)

/* Badges — top-level selectors, no deep nesting, for reliable scoping */
.badge
	display: inline-flex
	align-items: center
	gap: 6px
	padding: 4px 10px
	border-radius: 6px
	font-size: 12px
	font-weight: 500
	justify-content: center

	.badge-dot
		width: 6px
		height: 6px
		border-radius: 50%
		flex-shrink: 0

.badge-inactive
	background: rgba(248,113,113,0.12)
	border: 1px solid rgba(248,113,113,0.3)
	color: #f87171

	.badge-dot
		background: #f87171

.badge-active
	background: rgba(74,222,128,0.12)
	border: 1px solid rgba(74,222,128,0.25)
	color: #4ade80

	.badge-dot
		background: #4ade80
</style>
