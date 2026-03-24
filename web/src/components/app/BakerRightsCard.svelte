<script lang="ts">
	import { flip } from 'svelte/animate';
	import type { BlockRights } from '@src/common/types/status';
	import Card from '../starlight/components/Card.svelte';
	import { getBakerColor, normalizeBlockRights } from '@src/util/baker';

	export let mode: 'upcoming' | 'past';
	export let title = 'Baking Rights';
	export let rights: BlockRights[];
	export let showBakerColors = false;

	$: isUpcoming = mode === 'upcoming';
	$: cardAccent = isUpcoming
		? 'linear-gradient(90deg, var(--blue), var(--blue-dim))'
		: 'linear-gradient(90deg, var(--slate), var(--slate-dim))';
	$: cardHoverShadow = isUpcoming ? '0 8px 32px var(--blue-glow)' : '0 8px 32px var(--slate-glow)';
	$: cardColor = isUpcoming ? 'var(--blue)' : 'var(--slate)';
	$: cardIconBg = isUpcoming ? 'var(--blue-glow)' : 'var(--slate-glow)';
	$: cardIconBorder = isUpcoming ? 'rgba(96,165,250,0.3)' : 'rgba(148,163,184,0.2)';
</script>

<div
	class="rights-wrap"
	style:--card-accent={cardAccent}
	style:--card-hover-shadow={cardHoverShadow}
	style:--rights-color={cardColor}
	style:--rights-icon-bg={cardIconBg}
	style:--rights-icon-border={cardIconBorder}
>
	<Card>
		<div class="baker-rights">
			<div class="card-title">
				<span class="icon">
					{#if isUpcoming}
						<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
					{:else}
						<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
					{/if}
				</span>
				{title}
			</div>

			<div class="rights-header">
				<span>Block</span>
				<span>Validate</span>
				<span>Endorse</span>
			</div>

			<div class="block-rights-wrap">
				{#each rights as blockRights (blockRights.level)}
					<div animate:flip>
						{#each normalizeBlockRights(blockRights) as right}
							<div class="rights-row" class:upcoming={isUpcoming}>
								{#if showBakerColors}
									<div class="baker-dot" style:background={getBakerColor(right.baker)}></div>
								{/if}
								<span class="rights-block">{blockRights.level.toLocaleString()}</span>

								<div class="rights-stats">
									<span class="rights-count"
										class:no-rights={right.blocks === 0}
										class:warning={mode === 'past' && right.realizedBlocks === 0 && right.blocks > 0}
										class:success={mode === 'past' && right.realizedBlocks > 0}
									>
										{#if mode === 'past'}{right.blocks}/{right.realizedBlocks * right.blocks}{:else}{right.blocks}{/if}
									</span>
									<span class="rights-icon" class:dim={right.blocks === 0}>
										<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z"/></svg>
									</span>
								</div>

								<div class="rights-stats">
									<span class="rights-count"
										class:no-rights={right.attestations === 0}
										class:warning={mode === 'past' && right.realizedAttestations === 0 && right.attestations > 0}
										class:success={mode === 'past' && right.realizedAttestations > 0}
									>
										{#if mode === 'past'}{right.attestations}/{right.realizedAttestations * right.attestations}{:else}{right.attestations}{/if}
									</span>
									<span class="rights-icon" class:dim={right.attestations === 0}>
										<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z"/></svg>
									</span>
								</div>
							</div>
						{/each}
					</div>
				{/each}
			</div>
		</div>
	</Card>
</div>

<style lang="sass">
.rights-wrap
	height: 100%
	display: grid

	:global(.card)
		height: 100%
		box-sizing: border-box

.baker-rights
	display: grid
	grid-template-rows: auto auto 1fr
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
		color: var(--rights-color)

		.icon
			width: 20px
			height: 20px
			border-radius: 6px
			display: flex
			align-items: center
			justify-content: center
			background: var(--rights-icon-bg)
			border: 1px solid var(--rights-icon-border)
			flex-shrink: 0

	.rights-header
		display: grid
		grid-template-columns: 120px 1fr 1fr
		padding: 0 12px 8px
		border-bottom: 1px solid var(--border-default)

		span
			font-size: 10px
			font-weight: 600
			text-transform: uppercase
			letter-spacing: 1px
			color: var(--text-muted)

			&:not(:first-child)
				text-align: center

	.block-rights-wrap
		overflow-y: auto
		max-height: 320px

		&::-webkit-scrollbar
			width: 4px
		&::-webkit-scrollbar-track
			background: transparent
		&::-webkit-scrollbar-thumb
			background: var(--border-default)
			border-radius: 2px

		.rights-row
			display: grid
			grid-template-columns: 120px 1fr 1fr
			align-items: center
			padding: 10px 12px
			border-radius: 4px
			transition: background 0.15s
			font-size: 13px

			&:hover
				background: rgba(255,255,255,0.03)

			&.upcoming .rights-block
				color: var(--rights-color)

			.baker-dot
				width: 8px
				height: 8px
				border-radius: 20%
				margin-right: 6px

			.rights-block
				font-family: var(--font-mono)
				font-weight: 600
				color: var(--text-primary)

			.rights-stats
				display: flex
				align-items: center
				gap: 6px
				justify-content: center

				.rights-count
					font-family: var(--font-mono)
					font-weight: 500
					color: var(--text-secondary)

					&.no-rights
						opacity: 0.4

					&.warning
						color: var(--warning-color)

					&.success
						color: var(--success-color)

				.rights-icon
					width: 16px
					height: 16px
					display: flex
					align-items: center
					justify-content: center
					opacity: 0.5

					svg
						width: 16px
						height: 16px

					&.dim
						opacity: 0.25
</style>
