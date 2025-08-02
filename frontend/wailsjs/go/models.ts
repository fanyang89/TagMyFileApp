export namespace main {
	
	export class FileSystemItem {
	    key: string;
	    label: string;
	    path: string;
	    name: string;
	    isDirectory: boolean;
	    size: number;
	    modified: number;
	    isLeaf: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileSystemItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.path = source["path"];
	        this.name = source["name"];
	        this.isDirectory = source["isDirectory"];
	        this.size = source["size"];
	        this.modified = source["modified"];
	        this.isLeaf = source["isLeaf"];
	    }
	}

}

