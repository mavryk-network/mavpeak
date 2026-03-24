<script lang="ts">
	import { writeToClipboard } from '@src/util/clipboard';
	import Card from '@components/starlight/components/Card.svelte';
	import type { NodeStatus } from '@src/common/types/status';
	import { formatBlockHash, formatTimestampAgoStrict } from '@src/util/format';
	import { onDestroy } from 'svelte';

	export let node: NodeStatus;
	export let title = `Node`;

	$: blockTimestamp = formatTimestampAgoStrict(node.block?.timestamp ?? 0);
	$: isMvkt = title.toLowerCase() !== 'baker';
	$: cardColor = isMvkt ? 'var(--amber)' : 'var(--green)';
	$: cardAccent = isMvkt
		? 'linear-gradient(90deg, var(--amber), var(--amber-dim))'
		: 'linear-gradient(90deg, var(--green), var(--green-dim))';
	$: cardGlow = isMvkt ? 'var(--amber-glow)' : 'var(--green-glow)';
	$: cardIconBg = isMvkt ? 'rgba(251,191,36,0.15)' : 'rgba(74,222,128,0.15)';
	$: cardIconBorder = isMvkt ? 'rgba(251,191,36,0.3)' : 'rgba(74,222,128,0.3)';
	$: hoverShadow = isMvkt ? '0 8px 32px var(--amber-glow)' : '0 8px 32px var(--green-glow)';

	const interval = setInterval(() => {
		blockTimestamp = formatTimestampAgoStrict(node.block?.timestamp ?? 0);
	}, 500);

	onDestroy(() => clearInterval(interval));
</script>

<div
	class="node-wrap"
	style:--card-accent={cardAccent}
	style:--card-hover-shadow={hoverShadow}
	style:--node-color={cardColor}
	style:--node-icon-bg={cardIconBg}
	style:--node-icon-border={cardIconBorder}
>
	<Card>
		<div class="node-grid">
			<div class="card-title">
				<span class="icon">
					<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
				</span>
				{title}
				{#if isMvkt}
					<div class="online-badge" class:online={node.connection_status === 'connected'}>
						<div class="online-dot" class:connected={node.connection_status === 'connected'}></div>
						{node.connection_status === 'connected' ? 'Online' : node.connection_status}
					</div>
				{:else}
					<div class="conn-dot" class:connected={node.connection_status === 'connected'}></div>
				{/if}
			</div>

			<button class="unstyle-button address" on:click={() => writeToClipboard(node.address ?? '')}
				>{node.address}</button>

			{#if node.connection_status === 'connected'}
				<div class="chain-state">
					<div class="level">
						<span class="level-num">{(node.block?.level_info.level ?? 0).toLocaleString()}</span>
						<span class="cycle">#{node.block?.level_info.cycle}</span>
					</div>
					<button
						class="unstyle-button hash"
						on:click={() => writeToClipboard(node.block?.hash ?? '')}
					>{formatBlockHash(node.block?.hash ?? '')}</button>
					<div class="timestamp">{blockTimestamp}</div>
				</div>
			{:else}
				<div class="disconnected-status">DISCONNECTED</div>
			{/if}

			{#if node.network_info}
				<div class="network-info">
					<span class="connections-num">{node.network_info?.connection_count}</span>
					<span class="connections-label">connections</span>
					<div class="connections-bar">
						<div class="connections-bar-fill" style:width="{Math.min((node.network_info?.connection_count ?? 0) / 50 * 100, 100)}%"></div>
					</div>
				</div>
			{/if}
		</div>
	</Card>
</div>

<style lang="sass">
.node-wrap
	height: 100%
	display: grid

	:global(.card)
		height: 100%
		box-sizing: border-box

.node-grid
	display: grid
	grid-template-rows: auto auto 1fr auto
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
		color: var(--node-color)

		.icon
			width: 20px
			height: 20px
			border-radius: 6px
			display: flex
			align-items: center
			justify-content: center
			background: var(--node-icon-bg)
			border: 1px solid var(--node-icon-border)
			flex-shrink: 0

		.conn-dot
			margin-left: auto
			width: 8px
			height: 8px
			border-radius: 50%
			background: var(--error-color)
			animation: pulse-node 2s ease-in-out infinite

			&.connected
				background: var(--success-color)

		.online-badge
			margin-left: auto
			display: flex
			align-items: center
			gap: 6px
			font-size: 12px
			font-weight: 500
			color: var(--text-muted)

			&.online
				color: var(--success-color)

			.online-dot
				width: 6px
				height: 6px
				border-radius: 50%
				background: var(--error-color)
				animation: pulse-node 2s ease-in-out infinite

				&.connected
					background: var(--success-color)

	.address
		font-family: var(--font-mono)
		font-size: 12px
		color: var(--node-color)
		padding: 6px 10px
		background: rgba(0,0,0,0.2)
		border-radius: 6px
		border: 1px solid rgba(255,255,255,0.06)
		text-align: left

	.chain-state
		display: grid
		gap: 6px

		.level
			display: flex
			align-items: baseline
			gap: 10px

			.level-num
				font-family: var(--font-mono)
				font-size: 32px
				font-weight: 700
				color: var(--text-primary)

			.cycle
				font-family: var(--font-mono)
				font-size: 14px
				color: var(--node-color)
				padding: 3px 8px
				background: var(--node-icon-bg)
				border-radius: 6px
				border: 1px solid var(--node-icon-border)

		.hash
			font-family: var(--font-mono)
			font-size: 13px
			color: var(--text-secondary)
			text-align: left

		.timestamp
			font-size: 13px
			color: var(--text-muted)

	.disconnected-status
		display: flex
		align-items: center
		justify-content: center
		font-size: 1.25rem
		color: var(--error-color)
		height: 100%

	.network-info
		display: flex
		align-items: center
		gap: 10px
		padding: 10px 14px
		background: rgba(0,0,0,0.2)
		border: 1px solid var(--node-icon-border)
		border-radius: 4px

		.connections-num
			font-family: var(--font-mono)
			font-size: 24px
			font-weight: 700
			color: var(--node-color)

		.connections-label
			font-size: 13px
			color: var(--text-muted)

		.connections-bar
			flex: 1
			height: 4px
			background: var(--border-default)
			border-radius: 2px
			overflow: hidden
			margin-left: 8px

			.connections-bar-fill
				height: 100%
				background: linear-gradient(90deg, var(--green), var(--cyan))
				border-radius: 2px
				transition: width 0.3s ease

@keyframes pulse-node
	0%, 100%
		opacity: 1
	50%
		opacity: 0.6
</style>
