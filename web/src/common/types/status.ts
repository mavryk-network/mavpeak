export type VotingPeriodInfo = {
	position: number;
	remaining: number;
	voting_period: {
		index: number;
		kind: string;
		start_position: number;
	};
};

export type NodeStatus = {
	address?: string;
	connection_status: 'connected' | 'disconnected' | 'connecting' | 'paused';
	block?: {
		hash: string;
		timestamp: string;
		//fitness: string
		level_info: {
			level: number;
			level_position: number;
			cycle: number;
			cycle_position: number;
		};
		protocol: string;
		voting_period_info: VotingPeriodInfo;
	};
	network_info?: {
		connection_count: number;
		stats?: {
			total_sent: string;
			total_recv: string;
			current_inflow: number;
			current_outflow: number;
		} | null;
	} | null;
	is_essential: boolean;
};

export type NodesStatus = { [key: string]: NodeStatus };

export type BlockRights = {
	level: number;
	rights: { [key: string]: Array<number> };
	realized_checked: boolean;
};

export type RightsStatus = {
	level: number;
	rights: Array<BlockRights>;
	realized_checked?: boolean;
};

export type AmiServiceInfo = {
	status: string | 'running';
	started: string;
	formattedTimestamp?: string;
};

export type ApplicationServices = { [key: string]: AmiServiceInfo };

export type ServicesStatus = {
	timestamp: number;
	applications?: { [key: string]: ApplicationServices };
};

export type BakerStatus = {
	full_balance: string;
	staked_balance: string;
	liquid_balance: string;
};

export type BakersStatus = {
	level: number;
	bakers: {
		[key: string]: BakerStatus;
	};
};

export type LedgerWalletStatus = {
	app_version: string;
	authorized: boolean;
	ledger: string;
	ledger_status: string;
	pkh: string;
};

export type MavbakeStatus = {
	rights: RightsStatus;
	services: ServicesStatus;
	bakers: BakersStatus;
	wallets: { [key: string]: LedgerWalletStatus };
};

export type WalletStatus = {
	address: string;
	balance: number;
	level: string;
};

export type MavpayStatus = {
	services: ServicesStatus;
	wallet: WalletStatus;
};

export type PeakStatus = {
	id?: string;
	modules: {
		mavbake: MavbakeStatus | undefined;
		mavpay: MavpayStatus | undefined;
	};
	nodes: NodesStatus;
};

export type StatusUpdate = {
	kind: 'full' | 'diff';
	data: PeakStatus;
};

export type NormalizedBlockRights = {
	baker: string;
	blocks: number;
	attestations: number;
	realizedBlocks: number;
	realizedAttestations: number;
};
