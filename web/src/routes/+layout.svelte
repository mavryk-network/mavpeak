<script lang="ts">
	import '@src/nodespecific.ts';
	import '@src/styles/default.sass';
	import '@xterm/xterm/css/xterm.css';

	import { APP_CONNECTION_STATUS } from '@app/state/index';
</script>

<div class="layout-grid">
	<header>
		<div class="title-wrap">
			<a class="unstyle-link" href="/about">
				<img src="/assets/mavpeak-logo-white.png" alt="MavPeak" class="nav-logo" />
			</a>
			<div
				class="connection-status"
				class:connected={$APP_CONNECTION_STATUS === 'connected'}
				class:reconnecting={$APP_CONNECTION_STATUS === 'reconnecting'}
				class:paused={$APP_CONNECTION_STATUS === 'paused'}
			>
				<div class="connection-status-sign"></div>
				{$APP_CONNECTION_STATUS}
			</div>
		</div>
	</header>

	<main>
		<slot />
	</main>
	<footer>
		© {new Date().getFullYear()} Mavryk Dynamics
	</footer>
</div>

<style lang="sass">
:root
	--menu-gap: var(--spacing)

.layout-grid
	position: relative
	display: grid
	width: 100vw
	height: 100vh
	grid-template-columns: minmax(100px, 1fr)
	grid-template-rows: var(--header-height) 1fr
	grid-template-areas: "header" "main" "footer"
	color: var(--text-color)

	header
		position: fixed
		height: var(--header-height)
		width: 100vw
		display: flex
		justify-content: left
		align-items: center
		z-index: 100
		background: rgba(10, 14, 23, 0.8)
		backdrop-filter: blur(12px)
		border-bottom: 1px solid var(--border-default)

		.title-wrap
			display: grid
			width: 100vw
			grid-template-columns: auto 1fr auto
			align-items: center
			padding: 0 var(--spacing-x2)

			.nav-logo
				height: 28px
				width: auto
				display: block

		.connection-status
			grid-column: 3
			display: flex
			align-items: center
			gap: 8px
			padding: 6px 14px
			border-radius: 20px
			font-size: 13px
			font-weight: 500
			text-transform: uppercase
			letter-spacing: 1px
			background: rgba(248, 113, 113, 0.1)
			border: 1px solid rgba(248, 113, 113, 0.3)
			color: var(--red)

			&.connected
				background: var(--green-glow)
				border-color: rgba(74, 222, 128, 0.3)
				color: var(--green)

			&.reconnecting, &.paused
				background: rgba(251, 191, 36, 0.1)
				border-color: rgba(251, 191, 36, 0.3)
				color: var(--amber)

			.connection-status-sign
				display: inline-block
				width: 8px
				height: 8px
				border-radius: 50%
				background-color: var(--red)
				animation: pulse 2s ease-in-out infinite

			&.connected .connection-status-sign
				background-color: var(--success-color)

			&.reconnecting .connection-status-sign, &.paused .connection-status-sign
				background-color: var(--amber)

	main
		position: relative
		grid-area: main
		overflow-x: clip

	footer
		position: fixed
		display: flex
		width: 100%
		grid-area: footer
		justify-content: center
		color: var(--text-muted)
		font-size: 12px
		padding-bottom: var(--spacing-f2)
		bottom: 0
		z-index: -1

@keyframes pulse
	0%, 100%
		opacity: 1
	50%
		opacity: 0.5
</style>
