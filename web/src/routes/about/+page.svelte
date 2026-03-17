<svelte:options immutable />

<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '@components/starlight/components/Button.svelte';
	import { MAVPEAK_VERSION, MAVPEAK_CODENAME } from '@src/common/constants';
	import { DiscordIcon, TwitterIcon, TelegramIcon } from '@src/components/app/icons';
</script>

<div class="about-wrap">
	<div class="about">
		<img src="/assets/mavpeak-logo-peak-white.png" alt="MavPeak" class="title" />
		<img src="/assets/mavpeak-logo-name-white.png" alt="MavPeak" class="title" />
		<h5 class="version not-selectable">
			{MAVPEAK_VERSION}
			{MAVPEAK_CODENAME}
		</h5>
		<div class="social">
			<a class="svg-icon" href="https://discord.gg/2PnCr62YgF" target="_blank" rel="noreferrer">
				<DiscordIcon />
			</a>
			<a class="svg-icon" href="https://x.com/MavrykDynamics" target="_blank" rel="noreferrer">
				<TwitterIcon />
			</a>
			<a class="svg-icon" href="https://t.me/MavrykNetwork" target="_blank" rel="noreferrer">
				<TelegramIcon />
			</a>
		</div>
		<div class="navigation">
			<div class="home">
				<Button label="Home" on:click={() => goto('/')} />
			</div>
		</div>
	</div>
</div>

<style lang="sass">
	.about-wrap
		height: 100%
		margin: 0px !important
		min-height: 300px !important
		overflow: hidden
		--footer-height: 30px
		color: var(--text-color)
		
		.about
			display: flex
			flex-direction: column
			justify-content: center
			align-items: center
			height: calc(100% - var(--footer-height))
			margin:0 !important
			
			.logo
				width: 350px
				max-width: 40vw
				max-height: 20vh
			.title
				height: 40px
				width: auto
				margin-top: min(var(--spacing-x2), 3vh)
				margin-bottom: min(var(--spacing), 3vh)
			.version
				margin-top: 0
				margin-bottom: min(var(--spacing-x2), 3vh)
				font-size: 1.5rem
				font-weight: 500
			.social
				display: grid
				grid-template-columns: auto auto auto
				grid-gap: var(--spacing)

				a
					margin: 0 0.5rem
					color: var(--text-color)
					transition: transform 0.1s ease-in-out
					border-radius: var(--border-radius)
					text-decoration: none
					width: 50px
					
					&:hover
						transform: scale(1.1)

.navigation
	display: grid
	padding: var(--spacing-x3)
	grid-template-columns: 1fr 200px 1fr

	.home
		grid-column: 2
</style>
