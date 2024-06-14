import { IndividualSimUI } from '../../individual_sim_ui';
import { EquipmentSpec } from '../../proto/common';
import { Database } from '../../proto_utils/database';
import { Importer } from '../importers';
import { BulkTab } from '../individual_sim_ui/bulk_tab';

export class BulkGearJsonImporter extends Importer {
	private readonly simUI: IndividualSimUI<any>;
	private readonly bulkUI: BulkTab;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkUI: BulkTab) {
		super(parent, simUI, 'Bag Item Import', true);
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.descriptionElem.appendChild(
			<>
				<p>Import bag items from a JSON file, which can be created by the WowSimsExporter in-game AddOn.</p>
				<p>To import, upload the file or paste the text below, then click, 'Import'.</p>
			</>,
		);
	}

	async onImport(data: string) {
		try {
			const equipment = EquipmentSpec.fromJsonString(data, { ignoreUnknownFields: true });
			if (equipment?.items?.length > 0) {
				const db = await Database.loadLeftoversIfNecessary(equipment);
				const items = equipment.items.filter(spec => spec.id > 0 && db.lookupItemSpec(spec));
				if (items.length > 0) {
					this.bulkUI.addItems(items);
				}
			}
			this.close();
		} catch (e: any) {
			console.warn(e);
			alert(e.toString());
		}
	}
}
