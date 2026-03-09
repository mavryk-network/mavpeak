<script lang="ts">
	import { nodes } from '@app/state/index';
	import {
		state as mavbakeStatus,
		bakers as mavbakeBakers,
		services as mavbakeServices,
		wallets as mavbakeWallets,
		futureBakingRights,
		pastBakingRights,
		votingPeriodInfo
	} from '@app/state/mavbake';
	import { state as mavpayStatus } from '@app/state/mavpay';
	import NodeStatusCard from '@components/app/NodeStatusCard.svelte';
	import BakerStatusCard from '@components/app/BakerStatusCard.svelte';
	import BakerRightsCard from '@components/app/BakerRightsCard.svelte';
	import ServicesStatusCard from '@components/app/ServicesStatusCard.svelte';

	import GovernancePeriodCard from '@src/components/app/GovernancePeriodCard.svelte';
	import PayoutsCard from '@src/components/app/PayoutsCard.svelte';
	import LedgerStatusCard from '@src/components/app/LedgerStatusCard.svelte';

	$: showBakerColors = $mavbakeBakers.length > 1;
	$: expandedBakingRights = $mavbakeBakers.length > 1;
</script>

<div class="dashboard-grid-wrap">
	<div class="dashboard-grid">
		{#if $mavbakeStatus}
			{#each $mavbakeBakers as [baker, info]}
				<BakerStatusCard baker={baker ?? {}} status={info} showColor={showBakerColors} />
			{/each}
			{#each $mavbakeWallets as [walletId, info]}
				<LedgerStatusCard id={walletId.toUpperCase()} {info} />
			{/each}
		{/if}
		{#if $mavpayStatus}
			<PayoutsCard />
		{/if}
		{#if $mavbakeStatus}
			<GovernancePeriodCard votingPeriodInfo={$votingPeriodInfo} />
		{/if}
		{#if Object.keys($mavbakeServices.applications ?? {}).length > 0}
			<ServicesStatusCard title="Baker's Services" services={$mavbakeServices} />
		{/if}
		{#each $nodes as [node, info]}
			<NodeStatusCard node={info} title={node} />
		{/each}
		{#if $mavbakeStatus}
			<div class="baker-rights" class:expanded={expandedBakingRights}>
				<BakerRightsCard
					mode="upcoming"
					rights={$futureBakingRights}
					{showBakerColors}
					title="Upcoming Baking Rights"
				/>
			</div>
			<div class="baker-rights" class:expanded={expandedBakingRights}>
				<BakerRightsCard
					mode="past"
					rights={$pastBakingRights}
					{showBakerColors}
					title="Past Baking Rights"
				/>
			</div>
		{/if}
	</div>
</div>

<style lang="sass">
.dashboard-grid-wrap
	display: grid
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)
	gap: var(--spacing)

	.dashboard-grid
		display: grid
		grid-column: 2
		grid-template-columns: repeat(auto-fill, minmax(450px, 1fr))
		gap: var(--spacing)

		.baker-rights
			display: grid
			grid-template-rows: 1fr

			&.expanded
				grid-column: 1 / -1
</style>
