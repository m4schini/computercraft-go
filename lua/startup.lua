local addr = ""

function log(prefix, message)
    local div = " : "
    local time = os.date("%Z%X")
    write(time .. div .. prefix .. div .. message .. "\n")
end

function printLine()
    local w, h = term.getSize()
    for i = 1, w-1 do
        write("-")
    end
    print("-")
end

function setup()
    local config = {}
    if fs.exists(".config") then
        config = setupFromFile()
    else
        config = setupFromUser()
        local file = fs.open(".config", "w")
        file.write(textutils.serialiseJSON(config))
        file.close()
    end

    addr = "ws://" .. config["host"] .. "/api/ws"

    print("INITIALIZATION COMPLETE:")
    print(textutils.serialise(config))
    printLine()
end

function setupFromFile()
    local file = fs.open(".config", "r")
    local config = textutils.unserialiseJSON(file.readAll())
    file.close()
    return config
end

function setupFromUser()
    local config = {}

    print("HOST? ")
    write("> ")
    config["host"] = read()

    return config
end

function connect()
    ws, err = http.websocket(addr)
    if not ws then
        log("!!", err)
    end

    printLine()
    print("INCOMING MESSAGES: ->")
    print("OUTGOING MESSAGES: <-")
    printLine()

    while true do
        local message, isBinary = ws.receive()
        if not isBinary then
            log("->", message:gsub("\n", ""))
            t = textutils.unserialiseJSON(message)

            local f, err = loadstring(t.func)
            if not err then
                local result = f()
                local resultJson = textutils.serialiseJSON(result)

                log("<-", resultJson)
                ws.send(resultJson)
            else
                print(err)
                ws.send(textutils.serialiseJSON({
                    err = err
                }))
            end
        end
    end
end

setup()

local backoff = 1
while not pcall(connect) do
    log("! ", "trying to reconnect")
    sleep(backoff)
end