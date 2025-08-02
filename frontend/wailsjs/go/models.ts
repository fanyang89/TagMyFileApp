export namespace main {
	
	export class DiskSpaceInfo {
	    totalBytes: number;
	    freeBytes: number;
	    usedBytes: number;
	    totalSpaceGB: number;
	    freeSpaceGB: number;
	    usedSpaceGB: number;
	    usagePercent: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskSpaceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalBytes = source["totalBytes"];
	        this.freeBytes = source["freeBytes"];
	        this.usedBytes = source["usedBytes"];
	        this.totalSpaceGB = source["totalSpaceGB"];
	        this.freeSpaceGB = source["freeSpaceGB"];
	        this.usedSpaceGB = source["usedSpaceGB"];
	        this.usagePercent = source["usagePercent"];
	    }
	}
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

