function test(stream)
    local function test_filter(rec)
        return rec["bandwidth"] > 1000
    end

    local function test_map(rec)
--        local names = record.bin_names(rec)
--        local data = map()
--        for i,value in ipairs(names) do
--            data[value] = rec[value]
--        end
--        return data
        return rec
    end

    return stream : filter(test_filter) : map(test_map)
end