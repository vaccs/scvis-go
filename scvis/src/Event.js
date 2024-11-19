class Event {
    constructor(eventNum, lineNum, path) {
        this.eventNum = eventNum;
        this.lineNum = lineNum;
        this.path = path;
    }
}

class FileCloseEvent extends Event {
    constructor(eventNum, lineNum, path, fileDescriptor) {
        super(eventNum, lineNum, path);
        this.fileDescriptor = fileDescriptor;
    }
}

class FileOpenEvent extends Event {
    constructor(eventNum, lineNum, path, fileDescriptor, flags, mode, is_sym_link) {
        super(eventNum, lineNum, path);
        this.fileDescriptor = fileDescriptor;
        this.flags = flags;
        this.mode = mode;
        this.is_sym_link = is_sym_link;
    }
}

class FileReadEvent extends Event {
    constructor(eventNum, lineNum, path, fileDescriptor, valid) {
        super(eventNum, lineNum, path);
        this.fileDescriptor = fileDescriptor;
        this.valid = valid;
    }
}

class FileWriteEvent extends Event {
    constructor(eventNum, lineNum, path, fileDescriptor, valid) {
        super(eventNum, lineNum, path);
        this.fileDescriptor = fileDescriptor;
        this.valid = valid;
    }
}

class FreeEvent extends Event {
    constructor(eventNum, lineNum, path, address) {
        super(eventNum, lineNum, path);
        this.address = address;
    }
}

class MallocEvent extends Event {
    constructor(eventNum, lineNum, path, address, size) {
        super(eventNum, lineNum, path);
        this.address = address;
        this.size = size;
    }
}

 class FunctionInvocationEvent extends Event {
    constructor(eventNum, lineNum, path, invocationLineNum, invocationPath, functionName, address) {
        super(eventNum, lineNum, path);
        this.invocationLineNum = invocationLineNum;
        this.invocationPath = invocationPath;
        this.functionName = functionName;
        this.address = address;
    }
}

class LineChangeEvent extends Event {
    constructor(eventNum, lineNum, path) {
        super(eventNum, lineNum, path);
    }
}

class OutputEvent extends Event {
    constructor(eventNum, lineNum, path, output) {
        super(eventNum, lineNum, path);
        this.output = output;
    }
}

class RegisterEvent extends Event {
    constructor(eventNum, lineNum, path, registerName, value) {
        super(eventNum, lineNum, path);
        this.registerName = registerName;
        this.value = value;
    }
}   

class ReturnAddressEvent extends Event {
    constructor(eventNum, lineNum, path, dynamicLink, address, functionName) {
        super(eventNum, lineNum, path);
        this.dynamicLink = dynamicLink;
        this.address = address;
        this.functionName = functionName;
    }
}

class ReturnEvent extends Event {
    constructor(eventNum, lineNum, path) {
        super(eventNum, lineNum, path);
    }
}   

class VariableAccessEvent extends Event {
    constructor(eventNum, lineNum, path, scope, address, varType,  variableName, pointsTo, value) {
        super(eventNum, lineNum, path);
        this.variableName = variableName;
        this.address = address;
        this.value = value;
        this.scope = scope;
        this.varType = varType;
        this.pointsTo = pointsTo;
    }
}

export {
    Event,
    FileCloseEvent,
    FileOpenEvent,
    FileReadEvent,
    FileWriteEvent,
    FreeEvent,
    MallocEvent,
    FunctionInvocationEvent,
    LineChangeEvent,
    OutputEvent,
    RegisterEvent,
    ReturnAddressEvent,
    ReturnEvent,
    VariableAccessEvent
};