<script lang="ts">
	import type { BakerStatus } from '@src/common/types/status';

	export let status: BakerStatus | undefined = undefined;

	function formatMVRK(mumav: string | undefined): string {
		try {
			const val = Number(BigInt(mumav || '0')) / 1_000_000;
			return new Intl.NumberFormat('en-US', { maximumFractionDigits: 0 }).format(val);
		} catch {
			return '0';
		}
	}

	$: stakedRatioPct = (() => {
		try {
			const full = Number(status?.full_balance || '0');
			const staked = Number(status?.staked_balance || '0');
			return full > 0 ? Math.round((staked / full) * 100) : 0;
		} catch {
			return 0;
		}
	})();
</script>

<div class="balance-bar">
	<div class="bar-body">
		<div class="bar-title">
			<span class="icon">
				<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/></svg>
			</span>
			Validator Balance
		</div>
		{#if status}
			<div class="stat">
				<span class="stat-label">Total Balance</span>
				<span class="stat-value">{formatMVRK(status.full_balance)} <span class="currency">MVRK</span></span>
			</div>
			<div class="divider"></div>
			<div class="stat">
				<span class="stat-label">Staked</span>
				<span class="stat-value">{formatMVRK(status.staked_balance)} <span class="currency">MVRK</span></span>
			</div>
			<div class="divider"></div>
			<div class="stat">
				<span class="stat-label">Available</span>
				<span class="stat-value">{formatMVRK(status.liquid_balance)} <span class="currency">MVRK</span></span>
			</div>
			<div class="divider"></div>
			<div class="ratio-wrap">
				<span class="ratio-label">Staked ratio</span>
				<div class="ratio-bar">
					<div class="ratio-bar-fill" style:width="{stakedRatioPct}%"></div>
				</div>
				<span class="ratio-pct">{stakedRatioPct}%</span>
			</div>
		{:else}
			<span class="loading-text">Fetching balance data...</span>
		{/if}
	</div>
</div>

<style lang="sass">
.balance-bar
	position: relative
	background: #151c2c
	border-radius: 6px
	border: 1px solid #1e293b
	overflow: hidden
	transition: border-color 0.2s, box-shadow 0.2s

	&:hover
		border-color: rgba(45, 212, 191, 0.3)
		box-shadow: 0 4px 20px rgba(45, 212, 191, 0.08)

	&::before
		content: ''
		position: absolute
		top: 0
		left: 0
		right: 0
		height: 3px
		background: linear-gradient(90deg, var(--teal), var(--teal-dim))
		border-radius: 6px 6px 0 0
		z-index: 1

	.bar-body
		padding: var(--card-vertical-spacing) var(--card-horizontal-spacing)
		display: flex
		align-items: center
		gap: 32px

	.bar-title
		align-self: flex-start
		font-size: 11px
		font-weight: 600
		text-transform: uppercase
		letter-spacing: 1.5px
		color: var(--teal)
		display: flex
		align-items: center
		gap: 8px
		white-space: nowrap

		.icon
			width: 20px
			height: 20px
			border-radius: 4px
			display: flex
			align-items: center
			justify-content: center
			background: var(--teal-glow)
			border: 1px solid rgba(45, 212, 191, 0.3)
			color: var(--teal)

	.stat
		display: flex
		flex-direction: column
		gap: 2px

		.stat-label
			font-size: 10px
			font-weight: 600
			text-transform: uppercase
			letter-spacing: 1px
			color: var(--text-muted)

		.stat-value
			font-family: var(--font-mono)
			font-size: 20px
			font-weight: 700
			color: var(--text-primary)

		.currency
			font-size: 12px
			font-weight: 500
			color: var(--teal)
			margin-left: 4px

	.divider
		width: 1px
		height: 36px
		background: var(--border-default)
		flex-shrink: 0

	.ratio-wrap
		flex: 1
		display: flex
		align-items: center
		gap: 12px
		min-width: 160px

		.ratio-label
			font-size: 11px
			color: var(--text-muted)
			white-space: nowrap

		.ratio-bar
			flex: 1
			height: 6px
			background: var(--border-default)
			border-radius: 3px
			overflow: hidden
			min-width: 60px

			.ratio-bar-fill
				height: 100%
				background: linear-gradient(90deg, var(--teal), var(--cyan-dim))
				border-radius: 3px
				transition: width 0.3s ease

		.ratio-pct
			font-family: var(--font-mono)
			font-size: 13px
			font-weight: 600
			color: var(--teal)
			white-space: nowrap

.loading-text
	font-size: 13px
	color: var(--text-muted)
	font-style: italic
</style>
