<script lang="ts">
	import { goto } from '$app/navigation';
	import Card from '@components/starlight/components/Card.svelte';
	import type { VotingPeriodInfo } from '@src/common/types/status';
	import { getVotingPeriodTimeLeft } from '@src/util/gov';
	import { onDestroy } from 'svelte';

	export let votingPeriodInfo: VotingPeriodInfo | undefined;

	$: timeLeft = getVotingPeriodTimeLeft(votingPeriodInfo);

	const interval = setInterval(() => {
		timeLeft = getVotingPeriodTimeLeft(votingPeriodInfo);
	}, 500);

	onDestroy(() => clearInterval(interval));
</script>

<div class="governance-wrap">
	<Card class="governance-card">
		<div class="governance">
			<div class="card-title">
				<span class="icon">
					<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/></svg>
				</span>
				Governance
			</div>

			{#if votingPeriodInfo}
				<div class="period-info">
					<div class="period-main">
						<span class="kind">{votingPeriodInfo?.voting_period.kind}</span>
						<span class="period-label">period #{votingPeriodInfo.voting_period.index}</span>
					</div>
					<div class="timer">
						<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
						ends in {timeLeft}
					</div>
				</div>
			{:else}
				<div class="no-data">NO DATA</div>
			{/if}

			<button class="btn-open" on:click={() => goto('/governance')}>OPEN</button>
		</div>
	</Card>
</div>

<style lang="sass">
.governance-wrap
	display: grid
	height: 100%
	user-select: none

	:global(.governance-card)
		height: 100%
		box-sizing: border-box
		--card-accent: linear-gradient(90deg, var(--purple), var(--purple-dim))
		--card-hover-shadow: 0 8px 32px var(--purple-glow)

.governance
	display: grid
	grid-template-rows: auto 1fr auto
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
		color: var(--purple)

		.icon
			width: 20px
			height: 20px
			border-radius: 6px
			display: flex
			align-items: center
			justify-content: center
			background: var(--purple-glow)
			border: 1px solid rgba(168, 85, 247, 0.3)
			flex-shrink: 0

	.period-info
		display: flex
		flex-direction: column
		gap: 12px
		justify-content: center

		.period-main
			display: flex
			align-items: baseline
			gap: 10px
			flex-wrap: wrap

			.kind
				font-size: 28px
				font-weight: 700
				color: var(--purple)
				text-transform: capitalize

			.period-label
				font-family: var(--font-mono)
				font-size: 14px
				color: var(--text-secondary)

		.timer
			display: inline-flex
			align-items: center
			gap: 6px
			padding: 6px 16px
			border-radius: 8px
			background: var(--purple-glow)
			border: 1px solid rgba(168, 85, 247, 0.2)
			color: var(--purple)
			font-size: 13px
			font-weight: 500
			width: fit-content

	.no-data
		display: flex
		align-items: center
		justify-content: center
		font-size: 1.5rem
		color: var(--text-muted)

	.btn-open
		width: 100%
		padding: 12px
		border: 1px solid rgba(168, 85, 247, 0.3)
		border-radius: 4px
		background: var(--purple-glow)
		color: var(--purple)
		font-family: 'Inter', sans-serif
		font-size: 14px
		font-weight: 600
		cursor: pointer
		transition: all 0.2s
		letter-spacing: 0.5px

		&:hover
			background: rgba(168, 85, 247, 0.25)
			border-color: var(--purple)
</style>
