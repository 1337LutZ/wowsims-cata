import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { RogueOptions_PoisonImbue as Poison } from '../core/proto/rogue.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { RogueSpecs } from '../core/proto_utils/utils';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MainHandImbue = <SpecType extends RogueSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, Poison>({
		fieldName: 'mhImbue',
		numColumns: 1,
		values: [
			{ value: Poison.NoPoison, tooltip: 'No Main Hand Poison' },
			{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
			{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
			{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
		],
	});

export const OffHandImbue = <SpecType extends RogueSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, Poison>({
		fieldName: 'ohImbue',
		numColumns: 1,
		values: [
			{ value: Poison.NoPoison, tooltip: 'No Off Hand Poison' },
			{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
			{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
			{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
		],
	});

export const StartingOverkillDuration = <SpecType extends RogueSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'startingOverkillDuration',
		label: 'Starting Overkill duration',
		labelTooltip: 'Initial Overkill buff duration at the start of each iteration.',
		showWhen: (player: Player<SpecType>) => player.getTalents().overkill || player.getTalents().masterOfSubtlety > 0,
	});

export const VanishBreakTime = <SpecType extends RogueSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'vanishBreakTime',
		label: 'Vanish Break Time',
		labelTooltip: 'Time it takes to start attacking after casting Vanish.',
		extraCssClasses: ['experimental'],
		showWhen: (player: Player<SpecType>) => player.getTalents().overkill || player.getTalents().masterOfSubtlety > 0,
	});

export const AssumeBleedActive = <SpecType extends RogueSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'assumeBleedActive',
		label: 'Assume Bleed Always Active',
		labelTooltip: "Assume bleed always exists for 'Hunger for Blood' activation. Otherwise will only calculate based on own garrote/rupture.",
		extraCssClasses: ['within-raid-sim-hide'],
		showWhen: (player: Player<SpecType>) => player.getTalents().hungerForBlood,
	});

export const ApplyPoisonsManually = <SpecType extends RogueSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'applyPoisonsManually',
		label: 'Configure poisons manually',
		labelTooltip: 'Prevent automatic poison configuration that is based on equipped weapons.',
	});
