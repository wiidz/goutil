package regionHelper

// 省 update a_region set lv = 1 and merge_name = name where province_code is null
// 市 update a_region set lv = 2 and parent_code = province_code where province_code is not null and city_code is null
// 区 update a_region set lv = 3 and parent_code = city_code where city_code is not null

// update a_region set merge_name = name where province_code is null
// UPDATE a_region AS cityUpdate INNER JOIN a_region AS provinceData ON cityUpdate.province_code=provinceData.code SET cityUpdate.merge_name=concat(provinceData.name,",",cityUpdate.name)WHERE cityUpdate.province_code IS NOT NULL AND cityUpdate.city_code IS NULL
// UPDATE a_region AS areaUpdate INNER JOIN a_region AS cityData ON areaUpdate.city_code=cityData.code SET areaUpdate.merge_name=concat(cityData.merge_name,",",areaUpdate.name)WHERE areaUpdate.city_code IS NOT NULL
