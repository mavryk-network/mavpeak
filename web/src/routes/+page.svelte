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
	import BakerRightsCard from '@components/app/BakerRightsCard.svelte';
	import ServicesStatusCard from '@components/app/ServicesStatusCard.svelte';
	import GovernancePeriodCard from '@src/components/app/GovernancePeriodCard.svelte';
	import PayoutsCard from '@src/components/app/PayoutsCard.svelte';
	import LedgerStatusCard from '@src/components/app/LedgerStatusCard.svelte';
	import BakerStatusCard from '@components/app/BakerStatusCard.svelte';
	import { onMount } from 'svelte';
	import type { BakerStatus } from '@src/common/types/status';

	let initialBakerBalances: [string, BakerStatus][] = [];

	onMount(async () => {
		try {
			const res = await fetch('/api/mavbake/bakers/balances');
			if (res.ok) {
				const data: Record<string, BakerStatus> = await res.json();
				initialBakerBalances = Object.entries(data);
			}
		} catch {}
	});

	$: displayBakers = $mavbakeBakers.length > 0 ? $mavbakeBakers : initialBakerBalances;
	$: showBakerColors = $mavbakeBakers.length > 1;
	$: bakerNodes = $nodes.filter(([t]) => t.toLowerCase() === 'baker');
	$: mvktNodes = $nodes.filter(([t]) => t.toLowerCase() !== 'baker');
	$: hasServices = Object.keys($mavbakeServices.applications ?? {}).length > 0;
</script>

<div class="dashboard-wrap">
	{#if $mavbakeStatus}
		{#if displayBakers.length > 0}
			{#each displayBakers as [baker, info]}
				<BakerStatusCard status={info} />
			{/each}
		{:else}
			<BakerStatusCard />
		{/if}
		{#each $mavbakeWallets as [walletId, info]}
			<LedgerStatusCard id={walletId.toUpperCase()} {info} />
		{/each}
	{/if}

	<!-- Top row: 3 columns -->
	<div class="top-row">
		<!-- Col 1: Validator Services -->
		<div class="col">
			{#if hasServices}
				<ServicesStatusCard title="Validator Services" services={$mavbakeServices} />
			{:else if $mavpayStatus}
				<PayoutsCard />
			{/if}
		</div>

		<!-- Col 2: Baker / Validator nodes -->
		<div class="col">
			{#each bakerNodes as [node, info]}
				<NodeStatusCard node={info} title={node} />
			{/each}
		</div>

		<!-- Col 3: MVKT stacked above Governance -->
		<div class="stacked-col">
			{#each mvktNodes as [node, info]}
				<NodeStatusCard node={info} title={node} />
			{/each}
			{#if $mavbakeStatus}
				<GovernancePeriodCard votingPeriodInfo={$votingPeriodInfo} />
			{/if}
		</div>
	</div>

	<!-- Bottom row: 2-col validation rights -->
	{#if $mavbakeStatus}
		<div class="bottom-row">
			<BakerRightsCard
				mode="upcoming"
				rights={$futureBakingRights}
				{showBakerColors}
				title="Upcoming Validation Rights"
			/>
			<BakerRightsCard
				mode="past"
				rights={$pastBakingRights}
				{showBakerColors}
				title="Past Validation Rights"
			/>
		</div>
	{/if}
</div>

<style lang="sass">
.dashboard-wrap
	display: flex
	flex-direction: column
	gap: var(--spacing)
	padding: var(--spacing) var(--spacing-x2) var(--spacing-x3)

	.top-row
		display: grid
		grid-template-columns: 1fr 1fr 1fr
		gap: var(--spacing)
		align-items: stretch

		.col
			display: flex
			flex-direction: column
			gap: var(--spacing)

		.stacked-col
			display: flex
			flex-direction: column
			gap: var(--spacing)

	.bottom-row
		display: grid
		grid-template-columns: 1fr 1fr
		gap: var(--spacing)
</style>
