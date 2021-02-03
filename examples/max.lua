--[[ function returning the max between two numbers --]]
function max(num1, num2)
  if (num1 > num2) then
     result = num1;
  else
     result = num2;
  end
  return result;
end

biggest = max(3, 6)

fil.setresult(biggest)
